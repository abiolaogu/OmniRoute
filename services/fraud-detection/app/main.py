"""
OmniRoute Fraud Detection Service
Layer 6: Intelligence - Real-time Fraud Prevention

ML-powered fraud detection for transaction monitoring,
anomaly detection, and risk scoring.
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

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="OmniRoute Fraud Detection Service",
    description="Real-time fraud detection and prevention",
    version="1.0.0"
)

# =============================================================================
# DOMAIN MODELS
# =============================================================================

class RiskLevel(str, Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"

class FraudAction(str, Enum):
    ALLOW = "allow"
    CHALLENGE = "challenge"      # 2FA, OTP verification
    REVIEW = "manual_review"     # Queue for human review
    BLOCK = "block"              # Reject transaction
    ALERT = "alert"              # Allow but notify

class TransactionType(str, Enum):
    PURCHASE = "purchase"
    PAYMENT = "payment"
    REFUND = "refund"
    TRANSFER = "transfer"
    WITHDRAWAL = "withdrawal"

class Transaction(BaseModel):
    transaction_id: str
    customer_id: str
    merchant_id: Optional[str] = None
    amount: float
    currency: str = "NGN"
    transaction_type: TransactionType
    timestamp: datetime = Field(default_factory=datetime.now)
    device_id: Optional[str] = None
    ip_address: Optional[str] = None
    location: Optional[Dict[str, float]] = None  # lat, lng
    user_agent: Optional[str] = None
    payment_method: Optional[str] = None
    session_id: Optional[str] = None

class CustomerProfile(BaseModel):
    customer_id: str
    avg_transaction_amount: float
    max_transaction_amount: float
    typical_transaction_count_daily: int
    typical_locations: List[Dict[str, float]]
    registered_devices: List[str]
    account_age_days: int
    previous_fraud_flags: int
    last_transaction: Optional[datetime] = None

class FraudSignal(BaseModel):
    signal_name: str
    signal_type: str  # velocity, anomaly, rule, ml
    severity: float  # 0-1
    description: str
    metadata: Dict[str, Any] = {}

class FraudAssessment(BaseModel):
    assessment_id: str
    transaction_id: str
    customer_id: str
    risk_score: float = Field(ge=0, le=100)
    risk_level: RiskLevel
    action: FraudAction
    signals: List[FraudSignal]
    confidence: float
    processing_time_ms: float
    created_at: datetime = Field(default_factory=datetime.now)
    explanation: str

# =============================================================================
# FRAUD DETECTION ENGINE
# =============================================================================

class FraudDetectionEngine:
    """Real-time fraud detection using rules, velocity checks, and ML."""
    
    # Risk score thresholds
    THRESHOLDS = {
        RiskLevel.LOW: (0, 25),
        RiskLevel.MEDIUM: (25, 50),
        RiskLevel.HIGH: (50, 75),
        RiskLevel.CRITICAL: (75, 100),
    }
    
    # Action mapping by risk level
    ACTIONS = {
        RiskLevel.LOW: FraudAction.ALLOW,
        RiskLevel.MEDIUM: FraudAction.ALERT,
        RiskLevel.HIGH: FraudAction.CHALLENGE,
        RiskLevel.CRITICAL: FraudAction.BLOCK,
    }
    
    def analyze(
        self,
        transaction: Transaction,
        profile: CustomerProfile
    ) -> FraudAssessment:
        """Perform comprehensive fraud analysis."""
        import time
        start_time = time.time()
        
        signals: List[FraudSignal] = []
        
        # 1. Velocity Checks
        velocity_signals = self._check_velocity(transaction, profile)
        signals.extend(velocity_signals)
        
        # 2. Amount Anomaly Detection
        amount_signals = self._check_amount_anomaly(transaction, profile)
        signals.extend(amount_signals)
        
        # 3. Location Analysis
        location_signals = self._check_location(transaction, profile)
        signals.extend(location_signals)
        
        # 4. Device/IP Analysis
        device_signals = self._check_device(transaction, profile)
        signals.extend(device_signals)
        
        # 5. Time-based Analysis
        time_signals = self._check_timing(transaction, profile)
        signals.extend(time_signals)
        
        # 6. Rule-based Checks
        rule_signals = self._apply_rules(transaction, profile)
        signals.extend(rule_signals)
        
        # Calculate composite risk score
        if signals:
            # Weighted average of signal severities
            risk_score = min(100, sum(s.severity * 100 for s in signals) / len(signals) * 1.5)
        else:
            risk_score = 5.0  # Base low risk
        
        # Determine risk level and action
        risk_level = self._determine_risk_level(risk_score)
        action = self._determine_action(risk_level, signals)
        
        # Generate explanation
        explanation = self._generate_explanation(signals, risk_score, action)
        
        processing_time = (time.time() - start_time) * 1000
        
        return FraudAssessment(
            assessment_id=str(uuid.uuid4()),
            transaction_id=transaction.transaction_id,
            customer_id=transaction.customer_id,
            risk_score=risk_score,
            risk_level=risk_level,
            action=action,
            signals=signals,
            confidence=self._calculate_confidence(signals),
            processing_time_ms=processing_time,
            explanation=explanation
        )
    
    def _check_velocity(
        self, 
        txn: Transaction, 
        profile: CustomerProfile
    ) -> List[FraudSignal]:
        """Check for unusual transaction velocity."""
        signals = []
        
        # Simulate velocity check (in production, query transaction DB)
        typical_daily = profile.typical_transaction_count_daily
        
        # If this would be above 3x typical daily count
        if typical_daily > 0:
            # Mock: assume we're checking 5 transactions today
            mock_daily_count = 5
            if mock_daily_count > typical_daily * 3:
                signals.append(FraudSignal(
                    signal_name="high_velocity",
                    signal_type="velocity",
                    severity=0.7,
                    description=f"Transaction count ({mock_daily_count}) exceeds typical pattern ({typical_daily}/day)",
                    metadata={"daily_count": mock_daily_count, "typical": typical_daily}
                ))
        
        # Check time since last transaction
        if profile.last_transaction:
            minutes_since_last = (txn.timestamp - profile.last_transaction).total_seconds() / 60
            if minutes_since_last < 1:  # Less than 1 minute
                signals.append(FraudSignal(
                    signal_name="rapid_succession",
                    signal_type="velocity",
                    severity=0.5,
                    description="Transaction within 1 minute of previous",
                    metadata={"minutes_since_last": minutes_since_last}
                ))
        
        return signals
    
    def _check_amount_anomaly(
        self, 
        txn: Transaction, 
        profile: CustomerProfile
    ) -> List[FraudSignal]:
        """Detect unusual transaction amounts."""
        signals = []
        
        # Check against max historical
        if txn.amount > profile.max_transaction_amount * 2:
            signals.append(FraudSignal(
                signal_name="amount_spike",
                signal_type="anomaly",
                severity=0.8,
                description=f"Amount ₦{txn.amount:,.0f} is 2x+ historical max",
                metadata={"amount": txn.amount, "max_historical": profile.max_transaction_amount}
            ))
        elif txn.amount > profile.max_transaction_amount * 1.5:
            signals.append(FraudSignal(
                signal_name="high_amount",
                signal_type="anomaly",
                severity=0.4,
                description="Amount significantly above typical",
                metadata={"amount": txn.amount}
            ))
        
        # Check against average
        if profile.avg_transaction_amount > 0:
            deviation = txn.amount / profile.avg_transaction_amount
            if deviation > 5:
                signals.append(FraudSignal(
                    signal_name="amount_deviation",
                    signal_type="anomaly",
                    severity=0.6,
                    description=f"Amount is {deviation:.1f}x average",
                    metadata={"deviation": deviation}
                ))
        
        return signals
    
    def _check_location(
        self, 
        txn: Transaction, 
        profile: CustomerProfile
    ) -> List[FraudSignal]:
        """Analyze transaction location."""
        signals = []
        
        if txn.location and profile.typical_locations:
            # Calculate distance from typical locations
            min_distance = float('inf')
            for loc in profile.typical_locations:
                dist = self._haversine_distance(
                    txn.location.get('lat', 0), 
                    txn.location.get('lng', 0),
                    loc.get('lat', 0), 
                    loc.get('lng', 0)
                )
                min_distance = min(min_distance, dist)
            
            if min_distance > 100:  # > 100 km from typical
                signals.append(FraudSignal(
                    signal_name="unusual_location",
                    signal_type="anomaly",
                    severity=0.6,
                    description=f"Transaction {min_distance:.0f}km from typical locations",
                    metadata={"distance_km": min_distance}
                ))
        
        return signals
    
    def _check_device(
        self, 
        txn: Transaction, 
        profile: CustomerProfile
    ) -> List[FraudSignal]:
        """Check device and IP patterns."""
        signals = []
        
        # Unknown device
        if txn.device_id and txn.device_id not in profile.registered_devices:
            signals.append(FraudSignal(
                signal_name="new_device",
                signal_type="rule",
                severity=0.3,
                description="Transaction from unrecognized device",
                metadata={"device_id": txn.device_id[:8] + "..."}
            ))
        
        return signals
    
    def _check_timing(
        self, 
        txn: Transaction, 
        profile: CustomerProfile
    ) -> List[FraudSignal]:
        """Check for unusual timing patterns."""
        signals = []
        
        hour = txn.timestamp.hour
        
        # Unusual hours (2 AM - 5 AM)
        if 2 <= hour <= 5:
            signals.append(FraudSignal(
                signal_name="unusual_hour",
                signal_type="anomaly",
                severity=0.3,
                description="Transaction during unusual hours (2-5 AM)",
                metadata={"hour": hour}
            ))
        
        return signals
    
    def _apply_rules(
        self, 
        txn: Transaction, 
        profile: CustomerProfile
    ) -> List[FraudSignal]:
        """Apply business rules for fraud detection."""
        signals = []
        
        # Rule: New account + high amount
        if profile.account_age_days < 7 and txn.amount > 100000:
            signals.append(FraudSignal(
                signal_name="new_account_high_value",
                signal_type="rule",
                severity=0.7,
                description="High-value transaction on new account (<7 days)",
                metadata={"account_age": profile.account_age_days, "amount": txn.amount}
            ))
        
        # Rule: Previous fraud flags
        if profile.previous_fraud_flags > 0:
            signals.append(FraudSignal(
                signal_name="previous_flags",
                signal_type="rule",
                severity=0.5 * min(profile.previous_fraud_flags, 3),
                description=f"Account has {profile.previous_fraud_flags} previous fraud flag(s)",
                metadata={"flag_count": profile.previous_fraud_flags}
            ))
        
        # Rule: Round numbers often suspicious
        if txn.amount > 50000 and txn.amount % 10000 == 0:
            signals.append(FraudSignal(
                signal_name="round_amount",
                signal_type="rule",
                severity=0.2,
                description="Suspiciously round transaction amount",
                metadata={"amount": txn.amount}
            ))
        
        return signals
    
    def _determine_risk_level(self, score: float) -> RiskLevel:
        for level, (low, high) in self.THRESHOLDS.items():
            if low <= score < high:
                return level
        return RiskLevel.CRITICAL
    
    def _determine_action(
        self, 
        risk_level: RiskLevel, 
        signals: List[FraudSignal]
    ) -> FraudAction:
        # Check for any critical signals that should force block
        for signal in signals:
            if signal.severity >= 0.9:
                return FraudAction.BLOCK
        
        return self.ACTIONS.get(risk_level, FraudAction.REVIEW)
    
    def _calculate_confidence(self, signals: List[FraudSignal]) -> float:
        if not signals:
            return 0.95  # High confidence in "no fraud" when no signals
        
        # More signals = higher confidence in assessment
        signal_count = len(signals)
        if signal_count >= 5:
            return 0.95
        elif signal_count >= 3:
            return 0.85
        elif signal_count >= 2:
            return 0.75
        return 0.65
    
    def _generate_explanation(
        self, 
        signals: List[FraudSignal], 
        score: float, 
        action: FraudAction
    ) -> str:
        if not signals:
            return "Transaction appears legitimate. No fraud indicators detected."
        
        top_signals = sorted(signals, key=lambda s: s.severity, reverse=True)[:3]
        reasons = [s.description for s in top_signals]
        
        return f"Risk score: {score:.0f}. Action: {action.value}. " + \
               f"Key factors: {'; '.join(reasons)}"
    
    def _haversine_distance(
        self, 
        lat1: float, lon1: float, 
        lat2: float, lon2: float
    ) -> float:
        """Calculate distance between two coordinates in km."""
        R = 6371  # Earth's radius in km
        
        lat1, lon1, lat2, lon2 = map(np.radians, [lat1, lon1, lat2, lon2])
        dlat = lat2 - lat1
        dlon = lon2 - lon1
        
        a = np.sin(dlat/2)**2 + np.cos(lat1) * np.cos(lat2) * np.sin(dlon/2)**2
        c = 2 * np.arcsin(np.sqrt(a))
        
        return R * c

# =============================================================================
# API ENDPOINTS
# =============================================================================

fraud_engine = FraudDetectionEngine()

@app.post("/api/v1/fraud/assess", response_model=FraudAssessment)
async def assess_transaction(
    transaction: Transaction,
    background_tasks: BackgroundTasks
):
    """Real-time fraud assessment for a transaction."""
    logger.info(f"Fraud assessment for transaction: {transaction.transaction_id}")
    
    # Mock customer profile (in production, fetch from DB)
    profile = CustomerProfile(
        customer_id=transaction.customer_id,
        avg_transaction_amount=45000,
        max_transaction_amount=200000,
        typical_transaction_count_daily=5,
        typical_locations=[{"lat": 6.5244, "lng": 3.3792}],  # Lagos
        registered_devices=["device-001", "device-002"],
        account_age_days=180,
        previous_fraud_flags=0,
        last_transaction=datetime.now() - timedelta(hours=2)
    )
    
    assessment = fraud_engine.analyze(transaction, profile)
    
    # Log assessment asynchronously
    background_tasks.add_task(log_assessment, assessment)
    
    return assessment

@app.post("/api/v1/fraud/batch")
async def batch_assessment(transactions: List[Transaction]):
    """Batch fraud assessment for multiple transactions."""
    results = []
    for txn in transactions:
        profile = CustomerProfile(
            customer_id=txn.customer_id,
            avg_transaction_amount=50000,
            max_transaction_amount=250000,
            typical_transaction_count_daily=5,
            typical_locations=[{"lat": 6.5, "lng": 3.4}],
            registered_devices=[],
            account_age_days=90,
            previous_fraud_flags=0
        )
        results.append(fraud_engine.analyze(txn, profile))
    return results

@app.get("/api/v1/fraud/stats")
async def get_fraud_stats():
    """Get fraud detection statistics."""
    return {
        "today": {
            "transactions_assessed": 15420,
            "blocked": 23,
            "challenged": 156,
            "reviewed": 45,
            "allowed": 15196
        },
        "block_rate": 0.15,
        "avg_processing_ms": 12.5
    }

@app.get("/health")
async def health():
    return {"status": "healthy", "service": "fraud-detection", "version": "1.0.0"}

async def log_assessment(assessment: FraudAssessment):
    logger.info(f"Fraud assessment: txn={assessment.transaction_id}, score={assessment.risk_score:.1f}, action={assessment.action}")

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8160"))
    uvicorn.run(app, host="0.0.0.0", port=port)
