# app.py - Clean Flask app using the pickle model
import os
import numpy as np
import pandas as pd
import faiss
import pickle
from flask import Flask, request, jsonify
import google.generativeai as genai


# 🔐 Configure Gemini
genai.configure(api_key="AIzaSyCKbo8TaMTFtYmKk2tluJEe6kLFYnKml7Q")
model = genai.GenerativeModel("gemini-1.5-flash")

# 📦 Load the trained model
print("🔄 Loading model...")
try:
    BASE_DIR = os.path.dirname(os.path.abspath(__file__))
    pkl_path = os.path.join(BASE_DIR, 'pkl', 'real_estate_model.pkl')

    with open(pkl_path, 'rb') as f:
        model_data = pickle.load(f)
    
    preprocessor = model_data['preprocessor']
    feature_weights = model_data['feature_weights']
    index = model_data['faiss_index']
    contact_df = model_data['contact_df']
    numeric_cols = model_data['numeric_cols']
    bool_cols = model_data['bool_cols']
    
    print("✅ Model loaded successfully!")
    print(f"📊 Loaded {len(contact_df)} contacts")
    
except FileNotFoundError:
    print("❌ Model file not found! Please run 'python train_model.py' first.")
    exit(1)
except Exception as e:
    print(f"❌ Error loading model: {e}")
    exit(1)

# 🔥 Flask app
app = Flask(__name__)
top_k = 10
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

        # ✨ Generate comprehensive summary
        contacts_summary = ""
        for i, contact_data in enumerate(contacts_data, 1):
            contact = contact_data["contact_details"]
            score = contact_data["similarity_score_percent"]
            contacts_summary += f"""
    Contact {i} (Index: {contact_data['contact_index']}, Match: {score}%):
    - Budget: {contact.get('price', 'N/A')}
    - Area: {contact.get('area_surface', 'N/A')} sqm
    - Property Type: {contact.get('property_type', 'N/A')}
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

   

    if __name__ == "__main__":
        app.run(debug=True)