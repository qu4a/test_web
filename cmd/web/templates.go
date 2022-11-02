// обернуть динмамические данные в новую структуру, которая бедт действовать
// как едина структура хранения данных
package main

import (
	"html/template"
	"path/filepath"
	"test_web/pkg/models"
)

// Создаем тип templateData, который будет действовать как хранилище для
// любых динамических данных, которые нужно передать HTML-шаблонам.
// На данный момент он содержит только одно поле, но мы добавим в него другие
// по мере развития нашего приложения.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet //добавили поле Snippets в структуру
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//инициализируем новую карту, которая будет храниться в кэш
	cache := map[string]*template.Template{}
	// Используем функцию filepath.Glob, чтобы получить срез всех файловых путей с
	// расширением '.page.tmpl'. По сути, мы получим список всех файлов шаблонов для страниц
	// нашего веб-приложения.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	//перебираем файлы шаблона от каждой страницы
	for _, page := range pages {
		//извлечение конечное названия файла и присваивание его переменной name
		name := filepath.Base(page)
		//обрабатываем итерируемый файл
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Используем метод ParseGlob для добавления всех каркасных шаблонов.
		// В нашем случае это только файл base.layout.tmpl (основная структура шаблона).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Используем метод ParseGlob для добавления всех вспомогательных шаблонов.
		// В нашем случае это footer.partial.tmpl "подвал" нашего шаблона.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Добавляем полученный набор шаблонов в кэш, используя название страницы
		// (например, home.page.tmpl) в качестве ключа для нашей карты.
		cache[name] = ts
	}
	return cache, err
}
