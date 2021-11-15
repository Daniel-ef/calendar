package queries

var UserCreate = `INSERT INTO calendar.users
 (user_id, email, phone, first_name, last_name, day_start, day_end) 
 VALUES(:user_id, :email, :phone, :first_name, :last_name, :day_start, :day_end)
 RETURNING user_id;`

var UserFind = `SELECT user_id, email, phone, first_name, last_name, day_start, day_end
 FROM calendar.users WHERE `
