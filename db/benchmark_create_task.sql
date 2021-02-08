DELETE FROM TASK;
DELETE FROM TASK_USER WHERE status <> 'admin';
INSERT INTO TASK_USER (id, name, email,status) values ('00000000-0000-0000-0000-000000000000', 'Pat Smith', 'somethingelse@something.com','status');
