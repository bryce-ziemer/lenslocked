/* INNER JOIN (JOIN) - relationship must exist in both tables! */
SELECT * FROM users
    JOIN sessions ON users.id = sessions.user_id;


/* LEFT JOIN (JOIN) - Includes all items in table we are starting with (in this case users) 
And only records in the table we are joining that actually have  arelationship with it 
(everything in left, None-everything in right)*/
SELECT * FROM users
    LEFT JOIN sessions ON users.id = sessions.user_id;


/* RIGHT JOIN (JOIN) - Left table on returned if mapped to right table, left table always returned */
SELECT * FROM users
    RIGHT JOIN sessions ON users.id = sessions.user_id;


/* FULL OUTER JOIN (JOIN) - return records from all tables, even if do not have relationships! 
Links together if there is a relationship*/
SELECT * FROM users
   FULL OUTER JOIN sessions ON users.id = sessions.user_id;
