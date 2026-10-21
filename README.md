# Mini Product Management — Go (Gin) + ReactJS + MySQL

Ứng dụng quản lý sản phẩm đơn giản để thực hành xây dựng ứng dụng web full-stack.

## 1. Công nghệ sử dụng

- **Backend:** Go (Golang), Gin, REST API
- **Frontend:** ReactJS, Vite
- **Database:** MySQL
- **Thư viện bổ trợ:** `gin-contrib/cors`, `go-sql-driver/mysql`, `godotenv`

## 2. Chức năng

- Xem danh sách sản phẩm.
- Tìm kiếm sản phẩm theo tên.
- Lọc sản phẩm theo danh mục.
- Xem chi tiết sản phẩm theo slug.
- Thêm, cập nhật và xóa sản phẩm.
- Kiểm tra dữ liệu đầu vào trước khi thao tác với database.

## 3. Cấu trúc thư mục

```text
mini-product-management/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── .env.example
│   └── internal/
│       ├── handler/
│       │   └── product.go
│       ├── model/
│       │   └── product.go
│       └── validation/
│           └── validation.go
├── database/
│   └── schema.sql
├── frontend/
│   ├── package.json
│   ├── index.html
│   ├── vite.config.js
│   └── src/
│       ├── main.jsx
│       ├── App.jsx
│       └── style.css
├── .gitignore
└── README.md
```

## 4. Yêu cầu cài đặt

- Go 1.23 trở lên
- Node.js và npm
- MySQL Server 8 trở lên
- Git (nếu muốn clone repository)

Kiểm tra các công cụ:

```powershell
go version
node -v
npm -v
mysql --version
```

## 5. Cách khởi động project

### Bước 1: Tạo database

Mở MySQL Workbench hoặc MySQL CLI và chạy toàn bộ file `database/schema.sql`. File sẽ tạo database `mini_product_management`, bảng `products` và một số dữ liệu mẫu.

Có thể chạy bằng PowerShell nếu lệnh `mysql` đã có trong PATH:

```powershell
mysql -u root -p < database/schema.sql
```

Nhập mật khẩu MySQL khi được hỏi.

### Bước 2: Cấu hình và chạy backend Go + Gin

Mở **terminal thứ nhất** tại thư mục `backend`:

```powershell
cd backend
Copy-Item .env.example .env
```

Mở file `backend/.env`, cập nhật thông tin đăng nhập MySQL của bạn. Ví dụ:

```env
PORT=8080
DB_USER=root
DB_PASSWORD=mat_khau_mysql_cua_ban
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=mini_product_management
```

Thay `mat_khau_mysql_cua_ban` bằng mật khẩu thực tế. Nếu MySQL của bạn dùng tài khoản khác, hãy sửa `DB_USER` tương ứng.

Cài các thư viện Go và chạy backend:

```powershell
go mod tidy
go run .
```

Backend chạy tại `http://localhost:8080`. Giữ terminal này mở.

> Nếu bạn đang dùng Command Prompt thay vì PowerShell, có thể tạo file `.env` bằng cách sao chép `.env.example` thủ công hoặc dùng lệnh `copy .env.example .env`.

### Bước 3: Chạy frontend ReactJS

Mở **terminal thứ hai** tại thư mục gốc của project, sau đó chạy:

```powershell
cd frontend
npm install
npm run dev
```

Mở địa chỉ Vite hiển thị trong terminal, thường là `http://localhost:5173`.

### Những lần chạy sau

Sau khi đã cài dependency và tạo `.env`, mỗi lần chạy project:

1. Đảm bảo MySQL Server đang chạy.
2. Terminal thứ nhất: `cd backend` rồi chạy `go mod tidy` `go run .`.
3. Terminal thứ hai: `cd frontend` rồi chạy `npm install` `npm run dev`.
4. Mở `http://localhost:5173` trên trình duyệt.

## 6. API endpoints

Các endpoint sử dụng prefix `/api/v1`.

| Method   | Endpoint                                | Mô tả                     |
| -------- | --------------------------------------- | ------------------------- |
| `GET`    | `/api/v1/products`                      | Lấy danh sách sản phẩm    |
| `GET`    | `/api/v1/products?search=mouse`         | Tìm theo tên sản phẩm     |
| `GET`    | `/api/v1/products?category=electronics` | Lọc theo danh mục         |
| `GET`    | `/api/v1/products/:slug`                | Xem chi tiết theo slug    |
| `POST`   | `/api/v1/products`                      | Tạo sản phẩm              |
| `PUT`    | `/api/v1/products/:id`                  | Cập nhật sản phẩm theo ID |
| `DELETE` | `/api/v1/products/:id`                  | Xóa sản phẩm theo ID      |
| `GET`    | `/api/v1/categories`                    | Lấy danh sách danh mục    |

Tham số `limit` dùng để giới hạn số lượng kết quả, mặc định là `10` và tối đa `100`.

### Ví dụ JSON khi tạo/cập nhật sản phẩm

```json
{
  "name": "Wireless Mouse",
  "slug": "wireless-mouse",
  "description": "A wireless mouse",
  "price": 19.99,
  "category": "electronics"
}
```

Danh mục hiện được hỗ trợ: `electronics`, `books`, `clothing`, `home`.

## 7. Một số lỗi thường gặp

### `Access denied for user`

Kiểm tra `DB_USER` và `DB_PASSWORD` trong `backend/.env`. Đảm bảo bạn đang sửa đúng file `.env`, không phải chỉ sửa `.env.example`.

### `cannot connect to MySQL`

Kiểm tra MySQL Server đang chạy, cổng `DB_PORT` chính xác và database đã được tạo bằng `database/schema.sql`.

### Frontend không gọi được API

Đảm bảo backend đang chạy ở cổng `8080`, frontend đang chạy ở cổng Vite và cấu hình CORS cho phép địa chỉ frontend.

### Thiếu thư viện Go

Chạy `go mod tidy` trong thư mục `backend` khi máy có kết nối Internet để tải dependency.

## 8. Bảo mật khi đẩy lên GitHub

- Không commit `backend/.env` vì file này chứa mật khẩu database.
- Chỉ commit `backend/.env.example` với thông tin mẫu.
- Kiểm tra `git status` trước khi commit.

## Mục tiêu học tập

Dự án giúp thực hành Gin router và route group, path/query parameters, JSON request/response, HTTP status codes, validation, CRUD với MySQL, biến môi trường và kết nối ReactJS với REST API.
