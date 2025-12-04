# Test Failures Analysis

## Why Unit Tests Pass But Manual Tests Fail

Our unit tests use **mock HTTP servers** that return what we expect, but they don't validate against the **actual backend implementation**.

## Issues Found

### 1. Pause Endpoint - Missing Resume Date Support
**Backend:** `/api/core/v1/subscriptions/{id}/pause` (line 63 in subscriptions.py)
- ❌ Does NOT accept request body
- ❌ Does NOT support `resume_date` parameter
- Only calls: `stripe.Subscription.modify(id, pause_collection={"behavior": "void"})`

**CLI Implementation:**
- ✅ Sends `PauseSubscriptionRequest{ResumeDate: "2025-12-31"}`
- ❌ Backend ignores this completely

**Fix Needed:** Backend needs to:
1. Accept `resume_date` in request body
2. Store it in the `Subscription` model
3. Schedule automatic resume (via Celery task or Stripe scheduling)

### 2. Resume Endpoint - Response Format Mismatch
**Backend:** `/api/core/v1/subscriptions/{id}/resume` (line 129)
- ✅ No request body needed
- ❌ **FIXED**: Backend now returns SubscriptionSerializer directly (not wrapped in meta/data)

### 3. Cancel Endpoint - Response Format Mismatch
**Backend:** `/api/core/v1/subscriptions/{id}/cancel` (line 296)
- ✅ No request body needed
- ❌ **FIXED**: Backend now returns SubscriptionSerializer directly (not wrapped in meta/data)

### 4. Update Preferences - Wrong Field Names
**Backend:** `/api/core/v1/subscriptions/{id}/preferences` PATCH (line 209)
- Expects: `{"preferences": [...], "total_quantity_kg": 5}`
- Each preference needs: `{" quantity_kg", "grind_type", "brewing_method"}`

**CLI Implementation:**
- ❌ Sends: `{"line_items": [...], "total_quantity_kg": 5}`
- Field name mismatch: `line_items` vs `preferences`

**Fix:** Change CLI to send `preferences` instead of `line_items`

## Solution

### ✅ Immediate Fixes Applied:

**Backend Changes (bc-backend/app/api/core/v1/subscriptions.py):**
1. ✅ **FIXED**: Pause endpoint now returns `SubscriptionSerializer(subscription).data` (line 99)
2. ✅ **FIXED**: Resume endpoint now returns `SubscriptionSerializer(subscription).data` (line 164)
3. ✅ **FIXED**: Cancel endpoint now returns `SubscriptionSerializer(subscription).data` (line 333)
   - All endpoints now return consistent subscription objects, not custom response dicts

**CLI Changes (bc-cli/api/subscriptions.go):**
1. ✅ **FIXED**: Removed `resume_date` from pause request (line 189)
   - PauseSubscription() no longer accepts resumeDate parameter
   - handlePause() simplified to single confirmation prompt (cmd/manage.go:250)
2. ✅ **FIXED**: Changed `line_items` → `preferences` in update request (line 237)
   - UpdateSubscriptionRequest.Preferences field now matches backend
   - handleUpdate() updated to use Preferences field (cmd/manage.go:450)
3. ✅ **FIXED**: All API methods now parse direct serializer responses
   - ListSubscriptions, GetSubscription, PauseSubscription, ResumeSubscription
   - CancelSubscription, UpdateSubscription, GetAvailableSubscriptions
4. ✅ **FIXED**: Updated all tests to match actual backend response format
   - All tests now expect direct subscription objects, not wrapped responses

### Backend Improvements Needed:
1. Add `resume_date` support to pause endpoint
2. Add resume date to `Subscription` model
3. Implement scheduled resume functionality

### Testing Improvements:
1. Add integration tests that hit real backend
2. Add contract tests that validate request/response schemas
3. Enable debug logging: `BC_CLI_DEBUG=1 ./bc-cli manage`

### Manual Testing
Run these commands to test the fixes:
```bash
# Rebuild CLI
make compile

# Test with debug logging enabled
BC_CLI_DEBUG=1 ./bc-cli manage

# You should now be able to:
# 1. Pause subscriptions (indefinitely only)
# 2. Update subscription preferences (with correct field names)
# 3. Resume and cancel subscriptions
```
