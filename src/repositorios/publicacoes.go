package repositorios

import (
	"api/internal/domain/entities"
	"database/sql"
)

type Publicacoes struct {
	db *sql.DB
}

func NovoRepositorioDePublicacoes(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

func (repositorio Publicacoes) Criar(post entities.Post) (uint64, error) {
	statement, erro := repositorio.db.Prepare("insert into publicacoes (titulo, conteudo, autor_id) values (?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(post.Title, post.Content, post.AuthorID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

func (repositorio Publicacoes) BuscarPorID(publicacaoID uint64) (entities.Post, error) {
	linha, erro := repositorio.db.Query(`
		select p.*, u.nick from
		publicacoes p inner join usuarios u
		on u.id = p.autor_id where p.id = ?
	`, publicacaoID)
	if erro != nil {
		return entities.Post{}, erro
	}

	defer linha.Close()

	var post entities.Post

	if linha.Next() {
		if erro = linha.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); erro != nil {
			return entities.Post{}, erro
		}
	}

	return post, nil
}

func (repositorio Publicacoes) Buscar(usuarioID uint64) ([]entities.Post, error) {
	linhas, erro := repositorio.db.Query(`
		select distinct p.*, u.nick from publicacoes p
		inner join usuarios u on u.id = p.autor_id
		inner join seguidores s on p.autor_id = s.usuario_id 
		where u.id = ? or s.seguidor_id = ?
		order by 1 desc
	`, usuarioID, usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var posts []entities.Post

	for linhas.Next() {
		var post entities.Post
		if erro = linhas.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (repositorio Publicacoes) Atualizar(publicacaoID uint64, post entities.Post) error {
	statement, erro := repositorio.db.Prepare("update publicacoes set titulo = ?, conteudo = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(post.Title, post.Content, publicacaoID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio Publicacoes) Deletar(publicacaoID uint64) error {
	statement, erro := repositorio.db.Query(
		"delete from publicacoes where id = ?", publicacaoID,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	return nil
}

func (repositorio Publicacoes) BuscarPorUsuario(usuarioID uint64) ([]entities.Post, error) {
	linhas, erro := repositorio.db.Query(`
		select p.*, u.nick from publicacoes p
		join usuarios u on u.id = p.autor_id
		where p.autor_id = ?
	`, usuarioID)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var posts []entities.Post

	for linhas.Next() {
		var post entities.Post
		if erro = linhas.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (repositorio Publicacoes) Curtir(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare("update publicacoes set curtidas = curtidas + 1 where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio Publicacoes) Descurtir(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare(`
		update publicacoes set curtidas = 
		CASE 
			WHEN curtidas > 0 THEN curtidas - 1 
			ELSE 0 
		END 
		where id = ?
	 `)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}

	return nil
}
