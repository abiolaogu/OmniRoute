"""
OmniRoute Credit Scoring ML Service
Layer 6: Intelligence - Machine Learning Credit Assessment

This service provides real-time credit scoring for retailers and businesses
based on transaction history, payment behavior, and alternative data sources.
"""

from fastapi import FastAPI, HTTPException, BackgroundTasks
from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any
from datetime import datetime, timedelta
from enum import Enum
import numpy as np
import logging
import uuid
import os

# =============================================================================
# CONFIGURATION
# =============================================================================

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="OmniRoute Credit Scoring Service",
    description="ML-powered credit assessment for B2B commerce",
    version="1.0.0"
)

# =============================================================================
# DOMAIN MODELS
# =============================================================================

class RiskBand(str, Enum):
    PRIME = "prime"                    # 750+ score, lowest risk
    NEAR_PRIME = "near_prime"          # 680-749, low risk
    SUBPRIME = "subprime"              # 580-679, moderate risk
    DEEP_SUBPRIME = "deep_subprime"    # 500-579, high risk
    UNSCOREABLE = "unscoreable"        # <500 or insufficient data

class CreditDecision(str, Enum):
    APPROVED = "approved"
    CONDITIONAL = "conditional"
    DECLINED = "declined"
    MANUAL_REVIEW = "manual_review"

class CreditRequest(BaseModel):
    customer_id: str
    requested_amount: float
    purpose: str = "trade_credit"
    tenure_days: int = 30

class TransactionHistory(BaseModel):
    customer_id: str
    total_orders: int
    total_value: float
    avg_order_value: float
    first_order_date: datetime
    last_order_date: datetime
    orders_last_30_days: int
    orders_last_90_days: int

class PaymentBehavior(BaseModel):
    customer_id: str
    total_payments: int
    on_time_payments: int
    late_payments: int
    default_count: int
    avg_days_to_pay: float
    current_outstanding: float
    max_ever_outstanding: float

class BusinessProfile(BaseModel):
    customer_id: str
    business_type: str
    years_in_business: float
    employee_count: int
    has_physical_store: bool
    is_verified: bool
    verification_level: str  # none, basic, advanced, premium
    social_score: Optional[float] = None
    referral_count: int = 0

class AlternativeData(BaseModel):
    customer_id: str
    mobile_money_activity: Optional[float] = None  # Monthly transaction volume
    utility_payments_score: Optional[float] = None  # 0-100, payment consistency
    social_connections: int = 0
    app_engagement_score: Optional[float] = None  # 0-100, platform usage
    geolocation_stability: Optional[float] = None  # 0-100, location consistency

class CreditScoreResult(BaseModel):
    request_id: str
    customer_id: str
    credit_score: int = Field(ge=300, le=850)
    risk_band: RiskBand
    decision: CreditDecision
    approved_amount: float
    interest_rate: float
    tenure_days: int
    credit_limit: float
    score_factors: List[Dict[str, Any]]
    recommendations: List[str]
    created_at: datetime
    valid_until: datetime

class ScoreBreakdown(BaseModel):
    category: str
    weight: float
    raw_score: float
    weighted_score: float
    factors: List[str]

# =============================================================================
# CREDIT SCORING ENGINE
# =============================================================================

class CreditScoringEngine:
    """
    Multi-factor credit scoring engine using ML models and rule-based logic.
    """
    
    # Score weights by category
    WEIGHTS = {
        "payment_history": 0.35,      # Most important - 35%
        "credit_utilization": 0.20,   # How much credit used - 20%
        "business_stability": 0.15,   # Years in business - 15%
        "transaction_pattern": 0.15,  # Order frequency/value - 15%
        "alternative_data": 0.10,     # Mobile money, etc - 10%
        "relationship_length": 0.05,  # Time on platform - 5%
    }
    
    # Risk band thresholds
    RISK_BANDS = {
        RiskBand.PRIME: (750, 850),
        RiskBand.NEAR_PRIME: (680, 749),
        RiskBand.SUBPRIME: (580, 679),
        RiskBand.DEEP_SUBPRIME: (500, 579),
        RiskBand.UNSCOREABLE: (300, 499),
    }
    
    # Interest rates by risk band (annual)
    INTEREST_RATES = {
        RiskBand.PRIME: 0.12,           # 12% per annum
        RiskBand.NEAR_PRIME: 0.18,      # 18% per annum
        RiskBand.SUBPRIME: 0.24,        # 24% per annum
        RiskBand.DEEP_SUBPRIME: 0.36,   # 36% per annum
        RiskBand.UNSCOREABLE: 0.48,     # 48% per annum
    }
    
    # Credit limit multipliers (of avg monthly GMV)
    LIMIT_MULTIPLIERS = {
        RiskBand.PRIME: 3.0,
        RiskBand.NEAR_PRIME: 2.0,
        RiskBand.SUBPRIME: 1.0,
        RiskBand.DEEP_SUBPRIME: 0.5,
        RiskBand.UNSCOREABLE: 0.0,
    }

    def calculate_score(
        self,
        transactions: TransactionHistory,
        payments: PaymentBehavior,
        business: BusinessProfile,
        alternative: Optional[AlternativeData] = None
    ) -> tuple[int, List[ScoreBreakdown]]:
        """Calculate comprehensive credit score."""
        
        breakdowns = []
        
        # 1. Payment History Score (35%)
        payment_score, payment_factors = self._score_payment_history(payments)
        breakdowns.append(ScoreBreakdown(
            category="Payment History",
            weight=self.WEIGHTS["payment_history"],
            raw_score=payment_score,
            weighted_score=payment_score * self.WEIGHTS["payment_history"],
            factors=payment_factors
        ))
        
        # 2. Credit Utilization Score (20%)
        utilization_score, util_factors = self._score_credit_utilization(payments, transactions)
        breakdowns.append(ScoreBreakdown(
            category="Credit Utilization",
            weight=self.WEIGHTS["credit_utilization"],
            raw_score=utilization_score,
            weighted_score=utilization_score * self.WEIGHTS["credit_utilization"],
            factors=util_factors
        ))
        
        # 3. Business Stability Score (15%)
        stability_score, stab_factors = self._score_business_stability(business)
        breakdowns.append(ScoreBreakdown(
            category="Business Stability",
            weight=self.WEIGHTS["business_stability"],
            raw_score=stability_score,
            weighted_score=stability_score * self.WEIGHTS["business_stability"],
            factors=stab_factors
        ))
        
        # 4. Transaction Pattern Score (15%)
        pattern_score, pattern_factors = self._score_transaction_pattern(transactions)
        breakdowns.append(ScoreBreakdown(
            category="Transaction Pattern",
            weight=self.WEIGHTS["transaction_pattern"],
            raw_score=pattern_score,
            weighted_score=pattern_score * self.WEIGHTS["transaction_pattern"],
            factors=pattern_factors
        ))
        
        # 5. Alternative Data Score (10%)
        alt_score, alt_factors = self._score_alternative_data(alternative)
        breakdowns.append(ScoreBreakdown(
            category="Alternative Data",
            weight=self.WEIGHTS["alternative_data"],
            raw_score=alt_score,
            weighted_score=alt_score * self.WEIGHTS["alternative_data"],
            factors=alt_factors
        ))
        
        # 6. Relationship Length Score (5%)
        rel_score, rel_factors = self._score_relationship_length(transactions)
        breakdowns.append(ScoreBreakdown(
            category="Relationship Length",
            weight=self.WEIGHTS["relationship_length"],
            raw_score=rel_score,
            weighted_score=rel_score * self.WEIGHTS["relationship_length"],
            factors=rel_factors
        ))
        
        # Calculate final score
        total_weighted = sum(b.weighted_score for b in breakdowns)
        
        # Scale to FICO-like range (300-850)
        # Raw weighted max = 100, we need to scale to 300-850
        final_score = int(300 + (total_weighted / 100) * 550)
        final_score = max(300, min(850, final_score))
        
        return final_score, breakdowns

    def _score_payment_history(self, payments: PaymentBehavior) -> tuple[float, List[str]]:
        """Score based on payment behavior. Max 100 points."""
        factors = []
        score = 100.0
        
        if payments.total_payments == 0:
            factors.append("No payment history")
            return 50.0, factors
        
        # On-time payment ratio
        on_time_ratio = payments.on_time_payments / payments.total_payments
        if on_time_ratio >= 0.95:
            factors.append("Excellent payment record (95%+ on time)")
        elif on_time_ratio >= 0.85:
            factors.append("Good payment record (85%+ on time)")
            score -= 10
        elif on_time_ratio >= 0.70:
            factors.append("Fair payment record")
            score -= 25
        else:
            factors.append("Poor payment history")
            score -= 40
        
        # Defaults heavily penalized
        if payments.default_count > 0:
            penalty = min(30, payments.default_count * 15)
            score -= penalty
            factors.append(f"{payments.default_count} default(s) on record")
        
        # Average days to pay
        if payments.avg_days_to_pay <= 7:
            factors.append("Very prompt payer")
            score += 5
        elif payments.avg_days_to_pay <= 14:
            factors.append("Prompt payer")
        elif payments.avg_days_to_pay <= 30:
            factors.append("Average payment speed")
            score -= 5
        else:
            factors.append("Slow payment pattern")
            score -= 15
        
        return max(0, min(100, score)), factors

    def _score_credit_utilization(
        self, 
        payments: PaymentBehavior, 
        transactions: TransactionHistory
    ) -> tuple[float, List[str]]:
        """Score based on credit utilization. Max 100 points."""
        factors = []
        
        if payments.max_ever_outstanding == 0:
            factors.append("No credit utilization history")
            return 70.0, factors
        
        # Current utilization ratio
        if payments.current_outstanding > 0 and payments.max_ever_outstanding > 0:
            utilization = payments.current_outstanding / payments.max_ever_outstanding
        else:
            utilization = 0
        
        # Lower utilization is better
        if utilization <= 0.30:
            score = 100
            factors.append("Low credit utilization (<30%)")
        elif utilization <= 0.50:
            score = 85
            factors.append("Moderate credit utilization (30-50%)")
        elif utilization <= 0.70:
            score = 65
            factors.append("High credit utilization (50-70%)")
        elif utilization <= 0.90:
            score = 45
            factors.append("Very high utilization (70-90%)")
        else:
            score = 25
            factors.append("Maxed out credit (>90%)")
        
        return score, factors

    def _score_business_stability(self, business: BusinessProfile) -> tuple[float, List[str]]:
        """Score based on business stability. Max 100 points."""
        factors = []
        score = 50.0  # Base score
        
        # Years in business
        if business.years_in_business >= 5:
            score += 25
            factors.append("Established business (5+ years)")
        elif business.years_in_business >= 2:
            score += 15
            factors.append("Growing business (2-5 years)")
        elif business.years_in_business >= 1:
            score += 5
            factors.append("New business (1-2 years)")
        else:
            factors.append("Very new business (<1 year)")
        
        # Physical presence
        if business.has_physical_store:
            score += 10
            factors.append("Has physical store location")
        
        # Verification level
        verification_scores = {
            "premium": 15,
            "advanced": 10,
            "basic": 5,
            "none": 0
        }
        v_score = verification_scores.get(business.verification_level, 0)
        if v_score > 0:
            score += v_score
            factors.append(f"{business.verification_level.title()} verification")
        
        # Employee count (proxy for size)
        if business.employee_count >= 10:
            score += 5
            factors.append("Medium-sized operation")
        
        return min(100, score), factors

    def _score_transaction_pattern(self, transactions: TransactionHistory) -> tuple[float, List[str]]:
        """Score based on transaction patterns. Max 100 points."""
        factors = []
        score = 50.0
        
        # Order frequency (last 30 days)
        if transactions.orders_last_30_days >= 20:
            score += 25
            factors.append("Very active customer (20+ orders/month)")
        elif transactions.orders_last_30_days >= 10:
            score += 15
            factors.append("Active customer (10-20 orders/month)")
        elif transactions.orders_last_30_days >= 5:
            score += 5
            factors.append("Regular customer (5-10 orders/month)")
        else:
            factors.append("Occasional buyer (<5 orders/month)")
        
        # GMV trend (90 vs 30 day comparison)
        if transactions.orders_last_90_days > 0:
            monthly_rate_90 = transactions.orders_last_90_days / 3
            monthly_rate_30 = transactions.orders_last_30_days
            
            if monthly_rate_30 > monthly_rate_90 * 1.2:
                score += 15
                factors.append("Increasing order volume")
            elif monthly_rate_30 < monthly_rate_90 * 0.8:
                score -= 10
                factors.append("Declining order volume")
            else:
                score += 5
                factors.append("Stable order volume")
        
        # Average order value
        if transactions.avg_order_value >= 100000:
            score += 10
            factors.append("High-value orders")
        elif transactions.avg_order_value >= 50000:
            score += 5
            factors.append("Medium-value orders")
        
        return min(100, score), factors

    def _score_alternative_data(self, alt: Optional[AlternativeData]) -> tuple[float, List[str]]:
        """Score based on alternative data sources. Max 100 points."""
        factors = []
        
        if alt is None:
            factors.append("No alternative data available")
            return 50.0, factors  # Neutral score
        
        scores = []
        
        if alt.mobile_money_activity is not None:
            if alt.mobile_money_activity >= 500000:
                scores.append(90)
                factors.append("Strong mobile money activity")
            elif alt.mobile_money_activity >= 100000:
                scores.append(70)
                factors.append("Moderate mobile money activity")
            else:
                scores.append(50)
        
        if alt.utility_payments_score is not None:
            scores.append(alt.utility_payments_score)
            if alt.utility_payments_score >= 80:
                factors.append("Consistent utility payments")
        
        if alt.app_engagement_score is not None:
            scores.append(alt.app_engagement_score)
            if alt.app_engagement_score >= 80:
                factors.append("High platform engagement")
        
        if alt.social_connections > 10:
            scores.append(80)
            factors.append("Strong business network")
        
        if scores:
            return np.mean(scores), factors
        return 50.0, factors

    def _score_relationship_length(self, transactions: TransactionHistory) -> tuple[float, List[str]]:
        """Score based on relationship with platform. Max 100 points."""
        factors = []
        
        days_on_platform = (datetime.now() - transactions.first_order_date).days
        
        if days_on_platform >= 730:  # 2+ years
            factors.append("Long-term customer (2+ years)")
            return 100.0, factors
        elif days_on_platform >= 365:  # 1+ year
            factors.append("Established customer (1-2 years)")
            return 80.0, factors
        elif days_on_platform >= 180:  # 6+ months
            factors.append("Growing relationship (6-12 months)")
            return 60.0, factors
        elif days_on_platform >= 90:  # 3+ months
            factors.append("New relationship (3-6 months)")
            return 40.0, factors
        else:
            factors.append("Very new customer (<3 months)")
            return 25.0, factors

    def determine_risk_band(self, score: int) -> RiskBand:
        """Map credit score to risk band."""
        for band, (low, high) in self.RISK_BANDS.items():
            if low <= score <= high:
                return band
        return RiskBand.UNSCOREABLE

    def calculate_credit_limit(
        self, 
        risk_band: RiskBand, 
        avg_monthly_gmv: float,
        current_outstanding: float
    ) -> float:
        """Calculate approved credit limit."""
        multiplier = self.LIMIT_MULTIPLIERS.get(risk_band, 0)
        max_limit = avg_monthly_gmv * multiplier
        available = max(0, max_limit - current_outstanding)
        return available

    def make_decision(
        self,
        score: int,
        risk_band: RiskBand,
        requested_amount: float,
        available_limit: float
    ) -> tuple[CreditDecision, float]:
        """Make credit decision and determine approved amount."""
        
        if risk_band == RiskBand.UNSCOREABLE:
            return CreditDecision.DECLINED, 0
        
        if risk_band == RiskBand.DEEP_SUBPRIME:
            if requested_amount <= available_limit * 0.5:
                return CreditDecision.CONDITIONAL, min(requested_amount, available_limit * 0.5)
            return CreditDecision.MANUAL_REVIEW, 0
        
        if risk_band == RiskBand.SUBPRIME:
            if requested_amount <= available_limit * 0.75:
                return CreditDecision.APPROVED, min(requested_amount, available_limit * 0.75)
            return CreditDecision.CONDITIONAL, available_limit * 0.75
        
        # Near Prime and Prime
        if requested_amount <= available_limit:
            return CreditDecision.APPROVED, requested_amount
        return CreditDecision.CONDITIONAL, available_limit

# =============================================================================
# API ENDPOINTS
# =============================================================================

scoring_engine = CreditScoringEngine()

@app.post("/api/v1/credit/score", response_model=CreditScoreResult)
async def calculate_credit_score(
    request: CreditRequest,
    background_tasks: BackgroundTasks
):
    """
    Calculate credit score and make credit decision.
    """
    logger.info(f"Credit score request for customer: {request.customer_id}")
    
    # In production, these would come from actual data services
    transactions = TransactionHistory(
        customer_id=request.customer_id,
        total_orders=150,
        total_value=5000000,
        avg_order_value=33333,
        first_order_date=datetime.now() - timedelta(days=400),
        last_order_date=datetime.now() - timedelta(days=2),
        orders_last_30_days=12,
        orders_last_90_days=35
    )
    
    payments = PaymentBehavior(
        customer_id=request.customer_id,
        total_payments=140,
        on_time_payments=128,
        late_payments=12,
        default_count=0,
        avg_days_to_pay=8.5,
        current_outstanding=150000,
        max_ever_outstanding=500000
    )
    
    business = BusinessProfile(
        customer_id=request.customer_id,
        business_type="retail_store",
        years_in_business=3.5,
        employee_count=4,
        has_physical_store=True,
        is_verified=True,
        verification_level="advanced",
        referral_count=8
    )
    
    alternative = AlternativeData(
        customer_id=request.customer_id,
        mobile_money_activity=250000,
        utility_payments_score=85,
        app_engagement_score=72,
        social_connections=15
    )
    
    # Calculate score
    score, breakdowns = scoring_engine.calculate_score(
        transactions, payments, business, alternative
    )
    
    # Determine risk band
    risk_band = scoring_engine.determine_risk_band(score)
    
    # Calculate credit limit
    avg_monthly_gmv = transactions.total_value / 12
    credit_limit = scoring_engine.calculate_credit_limit(
        risk_band, avg_monthly_gmv, payments.current_outstanding
    )
    
    # Make decision
    decision, approved_amount = scoring_engine.make_decision(
        score, risk_band, request.requested_amount, credit_limit
    )
    
    # Get interest rate
    interest_rate = scoring_engine.INTEREST_RATES.get(risk_band, 0.48)
    
    # Build recommendations
    recommendations = []
    if score < 750:
        if payments.on_time_payments / payments.total_payments < 0.95:
            recommendations.append("Improve payment timeliness to boost score")
        if not business.is_verified:
            recommendations.append("Complete business verification for higher limits")
        if transactions.orders_last_30_days < 10:
            recommendations.append("Increase order frequency to demonstrate activity")
    
    result = CreditScoreResult(
        request_id=str(uuid.uuid4()),
        customer_id=request.customer_id,
        credit_score=score,
        risk_band=risk_band,
        decision=decision,
        approved_amount=approved_amount,
        interest_rate=interest_rate,
        tenure_days=request.tenure_days,
        credit_limit=credit_limit,
        score_factors=[
            {"category": b.category, "score": b.raw_score, "weight": b.weight, "factors": b.factors}
            for b in breakdowns
        ],
        recommendations=recommendations,
        created_at=datetime.now(),
        valid_until=datetime.now() + timedelta(days=30)
    )
    
    # Log decision asynchronously
    background_tasks.add_task(log_credit_decision, result)
    
    return result

@app.get("/api/v1/credit/score/{customer_id}")
async def get_latest_score(customer_id: str):
    """Get the most recent credit score for a customer."""
    # In production, fetch from database
    return {"message": f"Latest score for {customer_id}", "score": 720}

@app.get("/api/v1/credit/history/{customer_id}")
async def get_score_history(customer_id: str, limit: int = 10):
    """Get credit score history for a customer."""
    return {"customer_id": customer_id, "history": []}

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "credit-scoring", "version": "1.0.0"}

async def log_credit_decision(result: CreditScoreResult):
    """Log credit decision for audit and analytics."""
    logger.info(f"Credit decision: customer={result.customer_id}, score={result.credit_score}, decision={result.decision}")

# =============================================================================
# MAIN
# =============================================================================

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8150"))
    uvicorn.run(app, host="0.0.0.0", port=port)
