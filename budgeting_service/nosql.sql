select sum(amount) transactions  where   and user_id=$1 and
      $1 start_time between end_time



SELECT
    SUM(amount) AS total_income,
FROM
    transactions
WHERE
    user_id = $1
  AND type = 'income'
  AND date BETWEEN $2 AND $3



--     - _id (ObjectId)
--     - user_id (UUID, Users jadvaliga havola)
--     - category_id (ObjectId, Categories to'plamiga havola)
-- - amount (float)
-- - period (enum: daily, weekly, monthly, yearly)
-- - start_date (timestamp)
-- - end_date (timestamp)


SELECT
    category_id,
    SUM(amount) AS total_amount,
    period,
 start_time,  -- Earliest start_time in the group
      end_time       -- Latest end_time in the group
FROM
    budgets
WHERE
    user_id = $1
GROUP BY
    category_id,
    period





--                                           - _id (ObjectId)
--                                                    - user_id (UUID, Users jadvaliga havola)
--                                                          - name (string)
--                                                                     - target_amount (float)
--                                                             - current_amount (float)
--                                                        - deadline (timestamp)
--                                                         - status (enum: in_progress, achieved, failed)
--                                                      - created_at (timestamp)

SELECT
    SUM(target_amount) AS target_amount_sum,
    SUM(current_amount) AS current_amount_sum,
    SUM(target_amount) + SUM(current_amount) AS total_amount,
    status
FROM
    goals
WHERE
    user_id = $1
  AND deadline BETWEEN $2 AND $3
GROUP BY
    status;
