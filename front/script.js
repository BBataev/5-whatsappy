// script.js

// Сохраняем имя пользователя после успешного логина
function saveUsername(username) {
  localStorage.setItem('username', username)
}

// Получаем имя пользователя
function getUsername() {
  return localStorage.getItem('username')
}

async function checkAuth() {
  try {
    const res = await fetch('http://localhost:8080/api/me', {
      method: 'GET',
      credentials: 'include' 
    })

    if (!res.ok) {
      window.location.href = 'index.html'
      return null
    }

    const data = await res.json()
    return data
  } catch (err) {
    console.error('Auth check failed:', err)
    window.location.href = 'index.html'
    return null
  }
}
