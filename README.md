![CI](https://github.com/arosyri/architecture-lab-3/actions/workflows/go.yml/badge.svg)
# Painter — графічний HTTP-сервер на Go

Проєкт реалізує простий графічний сервер з веб-інтерфейсом, який приймає команди через HTTP та малює фігури і фон у вікні.

---

## Можливості

- Заливка фону (білий, зелений тощо)
- Малювання фігури T-180 жовтого кольору в заданій позиції
- Малювання прямокутника (bgrect)
- Малювання кольорової рамки (border)
- Переміщення фігури (move)
- Оновлення зображення (update)
- Скидання до початкового стану (reset)

---

## Встановлення та запуск

1. Клонуйте репозиторій:

```bash
git clone https://github.com/arosyri/architecture-lab-3.git
cd architecture-lab-3
```
2. Зберіть проєкт:

```bash
go build -o painter ./cmd/painter
```
3. Запустіть сервер:

```bash
./painter
```
## Використання
Замість того, щоб вручну вводити багато curl-запитів, можна зберегти команди у текстовий файл (cmd.txt) і надіслати їх через POST-запит.

Приклад cmd.txt

```
white
border green
bgrect 0.25 0.25 0.75 0.75
figure 400 400
move 500 500
figure 420 420
update
```
Запуск через curl (POST)

```bash
curl -X POST --data-binary @cmd.txt http://localhost:17000
```
## Тестування
Для запуску тестів виконайте:

```bash
go test ./... -v
```
## Діаграма залежностей
Файл .pdf у корені проєкту містить діаграму залежностей компонентів.
