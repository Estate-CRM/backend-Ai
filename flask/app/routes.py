
import numpy as np
from sklearn.preprocessing import OneHotEncoder, MinMaxScaler
from sklearn.compose import ColumnTransformer
import faiss
import google.generativeai as genai

import pandas as pd

import os
import json
# 🔐 Set Gemini API Key
genai.configure(api_key="AIzaSyBggEffzJCePRlwUwGOzQFYy8RZZA8QGbE")  # Replace with your actual Gemini key
model = genai.GenerativeModel("gemini-1.5-flash")

BASE_DIR = os.path.dirname(os.path.abspath(__file__))  # Points to flask/app
csv_path = os.path.join(BASE_DIR, "static", "synthetic_algerian_contacts.csv")
contact_df = pd.read_csv(csv_path)
property_df = pd.read_csv(csv_path)

# 🧹 Drop unneeded columns
drop_cols = ['id', 'client_id', 'agent_id', 'description', 'images', 'created_at']
contact_df.drop(columns=[col for col in drop_cols if col in contact_df.columns], inplace=True)
property_df.drop(columns=[col for col in drop_cols if col in property_df.columns], inplace=True)

# 🔗 Keep shared columns only
shared_columns = list(set(contact_df.columns) & set(property_df.columns))
contact_df = contact_df[shared_columns]
property_df = property_df[shared_columns]

# 🧠 Identify column types
categorical_cols = [col for col in shared_columns if contact_df[col].dtype == 'object']
bool_cols = [col for col in shared_columns if contact_df[col].dropna().isin([0, 1, True, False]).all() and col not in categorical_cols]
numeric_cols = [col for col in shared_columns if col not in categorical_cols + bool_cols]

# ⚖️ Feature importance weighting
important_features = [
    'price', 'area_surface',
    'latitude', 'longitude', 'property_type', 'has_parking',
    'distance_to_city_center', 'public_transport_accessible'
]
column_weights = {col: 1.0 if col in important_features else 0.95 for col in shared_columns}

# 🔧 Preprocessing pipeline
preprocessor = ColumnTransformer([
    ('cat', OneHotEncoder(handle_unknown='ignore'), categorical_cols),
    ('num', MinMaxScaler(), numeric_cols + bool_cols)
])

combined_df = pd.concat([contact_df, property_df], axis=0)
preprocessor.fit(combined_df)

X_contact_raw = preprocessor.transform(contact_df)
X_property_raw = preprocessor.transform(property_df.iloc[[0]])

# 🏷️ Get encoded feature names
encoder = preprocessor.named_transformers_['cat']
cat_feature_names = encoder.get_feature_names_out(categorical_cols)
num_feature_names = numeric_cols + bool_cols
all_feature_names = list(cat_feature_names) + num_feature_names

# 🎯 Apply weights to features
feature_weights = []
for name in all_feature_names:
    base_col = name.split('_')[0] if '_' in name else name
    weight = column_weights.get(base_col, 0.95)
    feature_weights.append(weight)
feature_weights = np.array(feature_weights, dtype=np.float32)

# ✅ Apply weights and convert to dense float32
X_contact = X_contact_raw * feature_weights
X_property = X_property_raw * feature_weights

X_contact = np.array(X_contact, dtype=np.float32)
X_property = np.array(X_property, dtype=np.float32)

# 🔄 Normalize for cosine similarity
faiss.normalize_L2(X_contact)
faiss.normalize_L2(X_property)

# 🔍 Create FAISS index
index = faiss.IndexFlatIP(X_contact.shape[1])
index.add(X_contact)

top_k = 10
scores, indices = index.search(X_property, top_k)

# 💬 Generate explanations with Gemini
print("\n📊 Top 10 Most Similar Contacts to the Property:\n")

for rank, (idx, score) in enumerate(zip(indices[0], scores[0]), 1):
    contact = contact_df.iloc[idx]
    prop = property_df.iloc[0]

    prompt = f"""
You are a real estate assistant AI. Compare the following contact preferences with a given property and explain why they are a good match.

CONTACT:
- Min Budget: {contact.get('min_budget', 'N/A')}
- Max Budget: {contact.get('max_budget', 'N/A')}
- Desired Area: {contact.get('desired_area_min', 'N/A')} to {contact.get('desired_area_max', 'N/A')}
- Location: ({contact.get('latitude', 'N/A')}, {contact.get('longitude', 'N/A')})
- Property Type: {contact.get('property_type', 'N/A')}
- Has Parking: {contact.get('has_parking', 'N/A')}
- Distance to City Center: {contact.get('distance_to_city_center', 'N/A')} km
- Public Transport Accessible: {contact.get('public_transport_accessible', 'N/A')}

PROPERTY:
- Price: {prop.get('price', 'N/A')}
- Area: {prop.get('area', 'N/A')}
- Location: ({prop.get('latitude', 'N/A')}, {prop.get('longitude', 'N/A')})
- Property Type: {prop.get('property_type', 'N/A')}
- Has Parking: {prop.get('has_parking', 'N/A')}
- Distance to City Center: {prop.get('distance_to_city_center', 'N/A')} km
- Public Transport Accessible: {prop.get('public_transport_accessible', 'N/A')}

💬 Explain in 2–3 sentences why this contact is similar to the property. Respond professionally.
    """

    try:
        response = model.generate_content(prompt)
        explanation = response.text.strip()
    except Exception as e:
        explanation = f"❌ Gemini error: {str(e)}"

    print(f"{rank}. Contact {idx} — Similarity: {score*100:.2f}%")
    print(f"   🧠 Explanation: {explanation}\n")

def setup_routes(app):
    @app.route("/predict", methods=["POST"])
    def predict():
        if request.is_json:
            user_input = request.get_json()
        else:
            user_input = {col: request.form.get(col) for col in shared_columns}

        user_df = pd.DataFrame([user_input])

        # Parse numeric and boolean types
        for col in numeric_cols + bool_cols:
            if col in user_df.columns:
                user_df[col] = pd.to_numeric(user_df[col], errors='coerce')

        # Transform and normalize
        X_property = preprocessor.transform(user_df).astype(np.float32)
        faiss.normalize_L2(X_property)

        # Perform similarity search
        scores, indices = index.search(X_property, 10)
        results = [
            {
                "contact_index": int(i),
                "similarity_score_percent": round(float(s) * 100, 2)
            }
            for i, s in zip(indices[0], scores[0])
        ]

        return json.dumps({
            "input": user_input,
            "results": results
        })