set database = tkdo;

DELETE FROM TASK;
DELETE FROM TASK_USER WHERE status <> 'admin';
INSERT INTO TASK_USER (id, name, email, client_secret, client_id, status) values ('00000000-0000-0000-0000-000000000000', 'Pat Smith', 'somethingelse@something.com','client_secret','client_id','status');
INSERT INTO TASK (id, name, type, status, user_id ) values ('60853a85-681d-4620-9677-946bbfdc8fbc', 'clean the gutters', 'basic|recurring', 'new', '00000000-0000-0000-0000-000000000000');
