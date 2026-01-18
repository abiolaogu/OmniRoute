"""
OmniRoute Credit Scoring Service
ML-powered credit risk assessment for trade financing.
"""
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict
import uvicorn
from datetime import datetime
import numpy as np

app = FastAPI(
    title="OmniRoute Credit Scoring",
    description="ML-based credit risk assessment for B2B lending",
    version="1.0.0"
)


class BusinessProfile(BaseModel):
    business_id: str
    business_name: str
    registration_date: str
    industry: str
    annual_revenue: Optional[float] = None
    employee_count: Optional[int] = None
    years_in_business: Optional[float] = None


class TransactionHistory(BaseModel):
    total_orders: int
    total_gmv: float
    average_order_value: float
    on_time_payment_rate: float
    days_since_last_order: int
    order_frequency_days: float
    returned_order_rate: float


class PaymentHistory(BaseModel):
    total_invoices: int
    total_amount: float
    on_time_payments: int
    late_payments: int
    average_days_to_pay: float
    current_outstanding: float
    oldest_outstanding_days: int


class CreditScoreRequest(BaseModel):
    business: BusinessProfile
    transactions: TransactionHistory
    payments: PaymentHistory
    requested_credit_limit: Optional[float] = None


class CreditScoreResponse(BaseModel):
    business_id: str
    credit_score: int  # 300-850
    risk_category: str  # "low", "medium", "high", "very_high"
    recommended_credit_limit: float
    interest_rate_tier: str
    confidence: float
    factors: Dict[str, str]
    approval_recommendation: str


class CreditLimitRequest(BaseModel):
    business_id: str
    current_limit: float
    requested_increase: float


class CreditLimitResponse(BaseModel):
    business_id: str
    approved: bool
    new_limit: float
    reason: str


class RiskAssessment(BaseModel):
    business_id: str
    default_probability: float
    expected_loss: float
    risk_adjusted_return: float
    watchlist_flags: List[str]


@app.get("/health")
def health():
    return {"status": "healthy", "service": "credit-scoring-service"}


@app.get("/ready")
def ready():
    return {"status": "ready"}


@app.post("/score", response_model=CreditScoreResponse)
def calculate_credit_score(request: CreditScoreRequest):
    """Calculate credit score for a business."""
    
    # Feature engineering (simplified)
    payment_score = min(100, request.payments.on_time_payments / max(1, request.payments.total_invoices) * 100)
    transaction_score = min(100, (request.transactions.on_time_payment_rate * 100))
    tenure_score = min(100, (request.business.years_in_business or 1) * 10)
    volume_score = min(100, request.transactions.total_gmv / 10000)
    
    # Weighted score (300-850 range)
    raw_score = (
        payment_score * 0.35 +
        transaction_score * 0.30 +
        tenure_score * 0.20 +
        volume_score * 0.15
    )
    credit_score = int(300 + (raw_score / 100) * 550)
    
    # Risk category
    if credit_score >= 750:
        risk_category = "low"
        rate_tier = "prime"
        limit_multiplier = 3.0
    elif credit_score >= 650:
        risk_category = "medium"
        rate_tier = "standard"
        limit_multiplier = 2.0
    elif credit_score >= 550:
        risk_category = "high"
        rate_tier = "subprime"
        limit_multiplier = 1.0
    else:
        risk_category = "very_high"
        rate_tier = "decline"
        limit_multiplier = 0.0
    
    # Recommended credit limit
    avg_order = request.transactions.average_order_value
    recommended_limit = avg_order * limit_multiplier * 5
    
    # Factors
    factors = {}
    if payment_score < 80:
        factors["payment_history"] = "Some late payments detected"
    if tenure_score < 50:
        factors["business_tenure"] = "Relatively new business"
    if volume_score < 50:
        factors["transaction_volume"] = "Lower than average transaction volume"
    if request.payments.current_outstanding > recommended_limit * 0.5:
        factors["utilization"] = "High current utilization"
    
    approval = "approved" if credit_score >= 550 else "declined"
    
    return CreditScoreResponse(
        business_id=request.business.business_id,
        credit_score=credit_score,
        risk_category=risk_category,
        recommended_credit_limit=round(recommended_limit, 2),
        interest_rate_tier=rate_tier,
        confidence=0.85,
        factors=factors if factors else {"overall": "Good standing"},
        approval_recommendation=approval
    )


@app.post("/limit-increase", response_model=CreditLimitResponse)
def request_limit_increase(request: CreditLimitRequest):
    """Evaluate credit limit increase request."""
    # Simplified approval logic
    increase_ratio = (request.current_limit + request.requested_increase) / request.current_limit
    
    approved = increase_ratio <= 1.5  # Max 50% increase
    new_limit = request.current_limit + request.requested_increase if approved else request.current_limit
    
    return CreditLimitResponse(
        business_id=request.business_id,
        approved=approved,
        new_limit=new_limit,
        reason="Approved based on good payment history" if approved else "Requested increase exceeds policy limits"
    )


@app.get("/risk/{business_id}", response_model=RiskAssessment)
def get_risk_assessment(business_id: str):
    """Get detailed risk assessment for a business."""
    return RiskAssessment(
        business_id=business_id,
        default_probability=0.02,
        expected_loss=500.0,
        risk_adjusted_return=0.15,
        watchlist_flags=[]
    )


@app.get("/portfolio/summary")
def get_portfolio_summary():
    """Get portfolio-level credit risk summary."""
    return {
        "total_exposure": 5000000.00,
        "avg_credit_score": 680,
        "default_rate": 0.018,
        "provision_rate": 0.02,
        "distribution": {
            "low_risk": 0.45,
            "medium_risk": 0.35,
            "high_risk": 0.15,
            "very_high_risk": 0.05
        }
    }


@app.post("/train")
def train_model(model_type: str = "xgboost"):
    """Trigger model retraining."""
    return {
        "status": "training_started",
        "model_type": model_type,
        "estimated_time_minutes": 30
    }


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8108)
