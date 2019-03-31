set database = tkdo;

DELETE FROM TASK;
DELETE FROM TASK_USER;
INSERT INTO TASK_USER (id, name, email) values ('00000000-0000-0000-0000-000000000000', 'Pat Smith', 'somethingelse@something.com');
INSERT INTO TASK (id, name, type, status, user_id ) values ('60853a85-681d-4620-9677-946bbfdc8fbc', 'clean the gutters', 'basic|recurring', 'new', '00000000-0000-0000-0000-000000000000');
