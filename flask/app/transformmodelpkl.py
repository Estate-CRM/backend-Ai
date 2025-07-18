# train_model.py - Run this once to create the pickle files
import os
import numpy as np
import pandas as pd
import faiss
import pickle
from sklearn.preprocessing import OneHotEncoder, MinMaxScaler
from sklearn.compose import ColumnTransformer

def train_and_save_model():
    print("🚀 Starting model training...")
    
    # 📁 Load and process data
    BASE_DIR = os.path.dirname(os.path.abspath(__file__))
    csv_path = os.path.join(BASE_DIR, "static", "synthetic_algerian_contacts.csv")

    contact_df = pd.read_csv(csv_path)
    property_df = pd.read_csv(csv_path)

    # Drop irrelevant columns
    drop_cols = ['id', 'client_id', 'agent_id', 'description', 'images', 'created_at']
    contact_df.drop(columns=[col for col in drop_cols if col in contact_df.columns], inplace=True)
    property_df.drop(columns=[col for col in drop_cols if col in property_df.columns], inplace=True)

    # Preprocess contact_df
    contact_df["price"] = (contact_df["min_budget"] + contact_df["max_budget"]) / 2
    contact_df.drop(columns=["min_budget", "max_budget"], inplace=True)
    contact_df["area_surface"] = (contact_df["desired_area_min"] + contact_df["desired_area_max"]) / 2
    contact_df.drop(columns=["desired_area_min", "desired_area_max"], inplace=True)

    # Shared columns
    shared_columns = list(set(contact_df.columns) & set(property_df.columns))
    contact_df = contact_df[shared_columns]
    property_df = property_df[shared_columns]

    # Type detection
    categorical_cols = [col for col in shared_columns if contact_df[col].dtype == 'object']
    bool_cols = [col for col in shared_columns if contact_df[col].dropna().isin([0, 1, True, False]).all() and col not in categorical_cols]
    numeric_cols = [col for col in shared_columns if col not in categorical_cols + bool_cols]

    important_features = [
        'price', 'area_surface', 'latitude', 'longitude',
        'property_type', 'has_parking', 'distance_to_city_center',
        'public_transport_accessible'
    ]
    column_weights = {col: 1.0 if col in important_features else 0.95 for col in shared_columns}

    # Preprocessing pipeline
    preprocessor = ColumnTransformer([
        ('cat', OneHotEncoder(handle_unknown='ignore'), categorical_cols),
        ('num', MinMaxScaler(), numeric_cols + bool_cols)
    ])
    combined_df = pd.concat([contact_df, property_df], axis=0)
    preprocessor.fit(combined_df)

    # Transform contact data
    X_contact_raw = preprocessor.transform(contact_df)

    # Feature weights
    encoder = preprocessor.named_transformers_['cat']
    cat_feature_names = encoder.get_feature_names_out(categorical_cols)
    num_feature_names = numeric_cols + bool_cols
    all_feature_names = list(cat_feature_names) + num_feature_names

    feature_weights = []
    for name in all_feature_names:
        base_col = name.split('_')[0] if '_' in name else name
        weight = column_weights.get(base_col, 0.95)
        feature_weights.append(weight)
    feature_weights = np.array(feature_weights, dtype=np.float32)

    # Apply weights
    X_contact = X_contact_raw * feature_weights
    X_contact = np.ascontiguousarray(X_contact.astype('float32'))

    # Normalize and create FAISS index
    faiss.normalize_L2(X_contact)
    index = faiss.IndexFlatIP(X_contact.shape[1])
    index.add(X_contact)

    # 💾 Save everything to pickle files
    model_data = {
        'preprocessor': preprocessor,
        'feature_weights': feature_weights,
        'faiss_index': index,
        'contact_df': contact_df,
        'numeric_cols': numeric_cols,
        'bool_cols': bool_cols,
        'categorical_cols': categorical_cols,
        'shared_columns': shared_columns
    }

    # Save to pickle file
    pkl_dir = os.path.join(BASE_DIR, "pkl")
    os.makedirs(pkl_dir, exist_ok=True)  # ✅ Ensure the folder exists

    model_path = os.path.join(pkl_dir, "real_estate_model.pkl")

    with open(model_path, 'wb') as f:
        pickle.dump(model_data, f)

    print("✅ Model saved to 'real_estate_model.pkl'")
    print(f"📊 Processed {len(contact_df)} contacts")
    print(f"🔍 FAISS index dimensions: {X_contact.shape}")

if __name__ == "__main__":
    train_and_save_model()