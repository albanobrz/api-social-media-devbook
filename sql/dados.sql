insert into users (name, nick, email, password)
values 
("usuário 1", "usuário1", "usuario1@gmail.com", "$2a$10$5kMk8gyWy50UDHYg5rOZlOVLRv663t/hYN/Afdwl1US9zlPPd6dAG"),
("usuário 2", "usuário2", "usuario2@gmail.com", "$2a$10$5kMk8gyWy50UDHYg5rOZlOVLRv663t/hYN/Afdwl1US9zlPPd6dAG"),
("usuário 3", "usuário3", "usuario3@gmail.com", "$2a$10$5kMk8gyWy50UDHYg5rOZlOVLRv663t/hYN/Afdwl1US9zlPPd6dAG");

insert into followers (user_id, follower_id)
values
(1, 2),
(3, 1),
(1, 3);

insert into posts(title, content, author_id)
values 
("Publicacao do usuario 1", "essa é a publicação do usuário 1, aha uhu", 1),
("Publicacao do usuario 2", "essa é a publicação do usuário 2, aha uhu", 2),
("Publicacao do usuario 3", "essa é a publicação do usuário 3, aha uhu", 3);
