import os
import numpy as np
import pandas as pd
import faiss
from flask import Flask, request, jsonify
from sklearn.preprocessing import OneHotEncoder, MinMaxScaler
from sklearn.compose import ColumnTransformer
import google.generativeai as genai

# 🔐 Configure Gemini
genai.configure(api_key="AIzaSyCKbo8TaMTFtYmKk2tluJEe6kLFYnKml7Q")  # Replace with your real API key
model = genai.GenerativeModel("gemini-1.5-flash")

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

# Normalize and index
faiss.normalize_L2(X_contact)
index = faiss.IndexFlatIP(X_contact.shape[1])
index.add(X_contact)
top_k = 10

# 🔥 Flask app
app = Flask(__name__)

def setup_routes(app):
    @app.route("/predict", methods=["POST"])
    def predict():
        if not request.is_json:
            return jsonify({"error": "Request must be JSON"}), 400

        user_input = request.get_json()
        user_df = pd.DataFrame([user_input])

        # Ensure numeric conversion
        for col in numeric_cols + bool_cols:
            if col in user_df.columns:
                user_df[col] = pd.to_numeric(user_df[col], errors='coerce')

        # Transform & weight user input
        X_user_raw = preprocessor.transform(user_df)
        X_user_weighted = X_user_raw * feature_weights
        X_user_weighted = np.ascontiguousarray(X_user_weighted.astype('float32'))
        faiss.normalize_L2(X_user_weighted)

        # Search top matches
        scores, indices = index.search(X_user_weighted, top_k)
        contacts_data = []

        for idx, score in zip(indices[0], scores[0]):
            if 0 <= idx < len(contact_df):
                contact = contact_df.iloc[idx].to_dict()
                contacts_data.append({
                    "contact_index": int(idx),
                    "similarity_score_percent": round(float(score) * 100, 2),
                    "contact_details": contact
                })

        # ✨ Generate one comprehensive summary for all contacts
        contacts_summary = ""
        for i, contact_data in enumerate(contacts_data, 1):
            contact = contact_data["contact_details"]
            score = contact_data["similarity_score_percent"]
            contacts_summary += f"""
Contact {i} (Index: {contact_data['contact_index']}, Match: {score}%):
- Budget: {contact.get('price', 'N/A')}
- Area: {contact.get('area_surface', 'N/A')} sqm
- Location: ({contact.get('latitude', 'N/A')}, {contact.get('longitude', 'N/A')})
- Property Type: {contact.get('property_type', 'N/A')}
- Has Parking: {contact.get('has_parking', 'N/A')}
- Distance to City Center: {contact.get('distance_to_city_center', 'N/A')} km
- Public Transport Accessible: {contact.get('public_transport_accessible', 'N/A')}
"""

        comprehensive_prompt = f"""
You are a real estate assistant AI. Analyze the following property and provide a brief summary of how well it matches with the top {len(contacts_data)} potential contacts.

PROPERTY DETAILS:
- Price: {user_input.get('price', 'N/A')}
- Area: {user_input.get('area_surface', 'N/A')} sqm
- Property Type: {user_input.get('property_type', 'N/A')}

TOP MATCHING CONTACTS:
{contacts_summary}

Provide ONLY two sentences: 
1. One sentence about the overall match quality
2. One sentence about the best contact match

Keep it brief and concise.
"""

        try:
            response = model.generate_content(comprehensive_prompt)
            comprehensive_analysis = response.text.strip()
        except Exception as e:
            comprehensive_analysis = f"❌ Gemini error: {str(e)}"

        return jsonify({
            "input": user_input,
            "total_matches": len(contacts_data),
            "contact_matches": contacts_data,
            "comprehensive_analysis": comprehensive_analysis
        })

# Mount the route to the app
setup_routes(app)

# Run it only if this file is the entry point
if __name__ == "__main__":
    app.run(debug=True)