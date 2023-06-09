SELECT profiles.id AS id,
       profiles.username AS username,
       profiles.slug AS slug,
       identities.first_name AS first_name,
       identities.last_name AS last_name,
       credentials.created_at AS created_at,
       COUNT(*) OVER() AS total
FROM credentials
    LEFT JOIN identities ON identities.id = credentials.id
    LEFT JOIN profiles ON profiles.id = credentials.id
    LEFT JOIN LATERAL (SELECT CASE WHEN ?0 = '' THEN '' ELSE format_user_search(?0) END AS term) AS parsed ON TRUE
    LEFT JOIN LATERAL (
        SELECT CASE
            WHEN parsed.term = '' THEN 1
            WHEN (profiles.username IS NOT NULL AND profiles.username <> '') THEN similarity(parsed.term, format_user_search(profiles.username))
            ELSE similarity(parsed.term, format_user_search(identities.first_name || ' ' || identities.last_name))
        END AS score
    ) AS username_proximity ON TRUE
    /* Slug has no accent or uppercase or special characters, so no need to format it */
    LEFT JOIN LATERAL (
        SELECT CASE
           WHEN parsed.term = '' THEN 1
           ELSE similarity(parsed.term, profiles.slug)
        END AS score
    ) AS slug_proximity ON TRUE
    LEFT JOIN LATERAL (SELECT GREATEST(username_proximity.score, slug_proximity.score) AS score) AS proximity ON TRUE
WHERE proximity.score > 0.1
ORDER BY proximity.score DESC, credentials.created_at DESC
LIMIT ?1 OFFSET ?2;
