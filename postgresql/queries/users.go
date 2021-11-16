package queries

var UserInsert = `INSERT INTO calendar.users
 (user_id, email, phone, first_name, last_name, workday_start, workday_end) 
 VALUES(:user_id, :email, :phone, :first_name, :last_name, :workday_start, :workday_end)
 RETURNING user_id;`

var UserSelect = `SELECT user_id, email, phone, first_name, last_name, workday_start, workday_end
 FROM calendar.users WHERE `
