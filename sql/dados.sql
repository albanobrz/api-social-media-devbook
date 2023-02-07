insert into usuarios (nome, nick, email, senha)
values 
("usuário 1", "usuário1", "usuario1@gmail.com", "$2a$10$5kMk8gyWy50UDHYg5rOZlOVLRv663t/hYN/Afdwl1US9zlPPd6dAG"),
("usuário 2", "usuário2", "usuario2@gmail.com", "$2a$10$5kMk8gyWy50UDHYg5rOZlOVLRv663t/hYN/Afdwl1US9zlPPd6dAG"),
("usuário 3", "usuário3", "usuario3@gmail.com", "$2a$10$5kMk8gyWy50UDHYg5rOZlOVLRv663t/hYN/Afdwl1US9zlPPd6dAG");

insert into seguidores (usuario_id, seguidor_id)
values
(1, 2),
(3, 1),
(1, 3);