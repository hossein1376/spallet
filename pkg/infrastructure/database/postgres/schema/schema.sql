CREATE TABLE users
(
    id       BIGSERIAL    NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE
);
CREATE TABLE transactions
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT      NOT NULL REFERENCES users (id) ON DELETE SET NULL,
    amount       NUMERIC     NOT NULL CHECK (amount > 0),
    description  TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    type         VARCHAR(31) NOT NULL,
-- only for deposit transactions
    release_date TIMESTAMPTZ,
-- only for withdrawal transactions
    status       VARCHAR(31),
    ref_id       UUID,
    CONSTRAINT chk_integrity CHECK (
        (type = 'deposit' AND status IS NULL AND ref_id IS NULL)
            OR
        (type = 'withdrawal' AND release_date IS NULL AND
         status IN ('pending', 'completed', 'failed'))
    )
);
CREATE TABLE balances
(
    user_id            BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    available          NUMERIC     NOT NULL DEFAULT 0 CHECK (available >= 0),
    total              NUMERIC     NOT NULL DEFAULT 0 CHECK (total >= 0),
    last_calculated_at TIMESTAMPTZ NOT NULL DEFAULT 'epoch'
);

CREATE INDEX idx_tx_user_created ON transactions (user_id, created_at DESC);
CREATE INDEX idx_tx_deposit ON transactions (user_id, release_date, created_at)
    WHERE type = 'deposit';
CREATE INDEX idx_tx_withdrawal ON transactions (user_id, created_at)
    WHERE type = 'withdrawal';
CREATE INDEX idx_balances_user ON balances (user_id);

CREATE FUNCTION refresh_user_balance(param_user_id BIGINT)
RETURNS TABLE (total_balance NUMERIC, available_balance NUMERIC) AS $$
DECLARE
    v_last TIMESTAMPTZ;
BEGIN
    SELECT last_calculated_at INTO v_last
    FROM balances WHERE user_id = param_user_id FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'user balance not found';
    END IF;

    UPDATE balances
    SET total = total + COALESCE((
        SELECT SUM(
            CASE
                WHEN type = 'deposit' THEN amount
                WHEN type = 'withdrawal' THEN -amount
                ELSE 0
            END
        )
        FROM transactions
        WHERE user_id = param_user_id AND created_at >= v_last
    ), 0),
    available = available + COALESCE((
        SELECT SUM(
            CASE
                WHEN type = 'deposit' AND (
                    (release_date IS NULL AND created_at >= v_last)
                        OR
                    (release_date >= v_last AND release_date <= NOW()))
                    THEN amount
                WHEN type = 'withdrawal' AND created_at >= v_last
                    THEN -amount
                ELSE 0
            END
        )
        FROM transactions
        WHERE user_id = param_user_id
    ), 0),
    last_calculated_at = NOW()
    WHERE user_id = param_user_id;

    RETURN QUERY
    SELECT total, available FROM balances WHERE user_id = param_user_id;
END;
$$ LANGUAGE plpgsql;

CREATE PROCEDURE refund_pending_transactions(p_tx_id BIGINT DEFAULT NULL)
LANGUAGE plpgsql
AS $$
DECLARE
    v_count INT;
BEGIN
    CREATE TEMPORARY TABLE pending_tx
    (
        id    BIGINT,
        user_id  BIGINT,
        amount   NUMERIC,
        ref_id   UUID
    ) ON COMMIT DROP;

    INSERT INTO pending_tx (id, user_id, amount, ref_id)
    SELECT id, user_id, amount, ref_id
    FROM transactions
    WHERE status = 'pending' AND (p_tx_id IS NULL OR id = p_tx_id)
    FOR UPDATE
    RETURNING 1 INTO v_count;

    IF p_tx_id IS NOT NULL AND v_count IS NULL THEN
        RAISE EXCEPTION 'Transaction not found or not pending';
    END IF;

    UPDATE transactions
    SET status = 'failed', updated_at = NOW()
    WHERE id IN (SELECT id FROM pending_tx);

    INSERT INTO transactions (user_id, amount, "type", description)
    SELECT user_id, amount, 'deposit', 'refund ' || ref_id::TEXT
    FROM pending_tx;
END;
$$;