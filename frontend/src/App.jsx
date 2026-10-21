import { useCallback, useEffect, useState } from 'react'

const API = 'http://localhost:8080/api/v1'
const emptyForm = { name: '', slug: '', description: '', price: '', category: 'electronics' }
const categories = [
  { value: 'electronics', label: 'Electronics' },
  { value: 'books', label: 'Books' },
  { value: 'clothing', label: 'Clothing' },
  { value: 'home', label: 'Home' },
]

async function request(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
  })
  if (response.status === 204) return null
  const body = await response.json().catch(() => ({}))
  if (!response.ok) throw new Error(body.error || 'Có lỗi xảy ra')
  return body
}

function App() {
  const [products, setProducts] = useState([])
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')
  const [form, setForm] = useState(emptyForm)
  const [editingId, setEditingId] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')

  const loadProducts = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const params = new URLSearchParams({ limit: '100' })
      if (search.trim()) params.set('search', search.trim())
      if (category) params.set('category', category)
      const result = await request(`${API}/products?${params.toString()}`)
      setProducts(result.data || [])
    } catch (err) {
      setError(`${err.message}. Hãy kiểm tra backend, MySQL và file .env.`)
    } finally {
      setLoading(false)
    }
  }, [search, category])

  useEffect(() => { loadProducts() }, [loadProducts])

  function updateField(event) {
    const { name, value } = event.target
    setForm(current => ({ ...current, [name]: value }))
  }

  function resetForm() {
    setForm(emptyForm)
    setEditingId(null)
  }

  async function submitForm(event) {
    event.preventDefault()
    setError('')
    setNotice('')
    const payload = { ...form, price: Number(form.price) }
    try {
      if (editingId) {
        await request(`${API}/products/${editingId}`, {
          method: 'PUT', body: JSON.stringify(payload),
        })
        setNotice('Đã cập nhật sản phẩm.')
      } else {
        await request(`${API}/products`, {
          method: 'POST', body: JSON.stringify(payload),
        })
        setNotice('Đã thêm sản phẩm.')
      }
      resetForm()
      await loadProducts()
    } catch (err) {
      setError(err.message)
    }
  }

  function startEdit(product) {
    setEditingId(product.id)
    setForm({
      name: product.name,
      slug: product.slug,
      description: product.description || '',
      price: String(product.price),
      category: product.category,
    })
    setNotice('')
    setError('')
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  async function deleteProduct(product) {
    if (!window.confirm(`Bạn có chắc muốn xóa "${product.name}"?`)) return
    setError('')
    setNotice('')
    try {
      await request(`${API}/products/${product.id}`, { method: 'DELETE' })
      setNotice('Đã xóa sản phẩm.')
      if (editingId === product.id) resetForm()
      await loadProducts()
    } catch (err) {
      setError(err.message)
    }
  }

  function makeSlug(value) {
    return value.normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '')
  }

  return (
    <main className="page">
      <header className="hero">
        <div>
          <p className="eyebrow">GO + GIN · REACT · MYSQL</p>
          <h1>Quản lý sản phẩm</h1>
          <p className="subtitle">Project luyện REST API, validation và CRUD với MySQL.</p>
        </div>
        <div className="count-card">
          <span>Tổng sản phẩm đang hiển thị</span>
          <strong>{products.length}</strong>
        </div>
      </header>

      {error && <div className="alert error">{error}</div>}
      {notice && <div className="alert success">{notice}</div>}

      <section className="panel">
        <div className="section-heading">
          <div>
            <p className="eyebrow">{editingId ? 'UPDATE' : 'CREATE'}</p>
            <h2>{editingId ? `Sửa sản phẩm #${editingId}` : 'Thêm sản phẩm mới'}</h2>
          </div>
          {editingId && <button className="button secondary" onClick={resetForm} type="button">Hủy sửa</button>}
        </div>

        <form onSubmit={submitForm} className="product-form">
          <label>
            Tên sản phẩm *
            <input name="name" value={form.name} onChange={updateField} minLength="2" maxLength="100" required placeholder="Ví dụ: Wireless Mouse" />
          </label>
          <label>
            Slug *
            <div className="slug-row">
              <input name="slug" value={form.slug} onChange={updateField} maxLength="120" required placeholder="wireless-mouse" />
              <button className="button secondary" type="button" onClick={() => setForm(current => ({ ...current, slug: makeSlug(current.name) }))}>Tạo slug</button>
            </div>
          </label>
          <label>
            Giá (USD) *
            <input name="price" type="number" min="0.01" step="0.01" value={form.price} onChange={updateField} required placeholder="19.99" />
          </label>
          <label>
            Danh mục *
            <select name="category" value={form.category} onChange={updateField} required>
              {categories.map(item => <option key={item.value} value={item.value}>{item.label}</option>)}
            </select>
          </label>
          <label className="full">
            Mô tả
            <textarea name="description" value={form.description} onChange={updateField} rows="3" placeholder="Mô tả ngắn về sản phẩm..." />
          </label>
          <div className="form-actions full">
            <button className="button primary" type="submit">{editingId ? 'Lưu thay đổi' : 'Thêm sản phẩm'}</button>
            {!editingId && <button className="button secondary" type="button" onClick={resetForm}>Làm mới</button>}
          </div>
        </form>
      </section>

      <section className="panel">
        <div className="section-heading">
          <div>
            <p className="eyebrow">PRODUCT INVENTORY</p>
            <h2>Danh sách sản phẩm</h2>
          </div>
          <button className="button secondary" onClick={loadProducts} type="button">Làm mới danh sách</button>
        </div>
        <div className="filters">
          <label className="search-box">
            Tìm theo tên
            <input value={search} onChange={event => setSearch(event.target.value)} placeholder="Nhập tên sản phẩm..." />
          </label>
          <label>
            Lọc danh mục
            <select value={category} onChange={event => setCategory(event.target.value)}>
              <option value="">Tất cả danh mục</option>
              {categories.map(item => <option key={item.value} value={item.value}>{item.label}</option>)}
            </select>
          </label>
        </div>

        {loading ? <p className="empty">Đang tải dữ liệu...</p> : products.length === 0 ? (
          <p className="empty">Chưa có sản phẩm phù hợp. Hãy thêm sản phẩm hoặc đổi từ khóa tìm kiếm.</p>
        ) : (
          <div className="table-wrap">
            <table>
              <thead><tr><th>Sản phẩm</th><th>Danh mục</th><th>Giá</th><th>Slug</th><th>Thao tác</th></tr></thead>
              <tbody>
                {products.map(product => (
                  <tr key={product.id}>
                    <td>
                      <div className="product-name">{product.name}</div>
                      <div className="description">{product.description || 'Chưa có mô tả'}</div>
                    </td>
                    <td><span className="category-tag">{product.category}</span></td>
                    <td className="price">${Number(product.price).toFixed(2)}</td>
                    <td><code>{product.slug}</code></td>
                    <td>
                      <div className="actions">
                        <button className="button small secondary" onClick={() => startEdit(product)} type="button">Sửa</button>
                        <button className="button small danger" onClick={() => deleteProduct(product)} type="button">Xóa</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
      <footer>Mini Product Management · Dùng để luyện tập full-stack cơ bản</footer>
    </main>
  )
}

export default App
