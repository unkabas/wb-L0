import { useState } from 'react'
import './App.css'

function App() {
	const [orderId, setOrderId] = useState('')
	const [order, setOrder] = useState(null)
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState('')

	const handleInputChange = e => {
		setOrderId(e.target.value)
	}

	const handleSearch = async () => {
		setLoading(true)
		setError('')
		setOrder(null)
		try {
			const res = await fetch(`http://localhost:8080/order/${orderId}`)
			if (!res.ok) throw new Error('Заказ не найден')
			const data = await res.json()
			setOrder(data)
		} catch (err) {
			setError(err.message)
		} finally {
			setLoading(false)
		}
	}

	return (
		<>
			<div>
				<h1 className='text-3xl font-extrabold'>Поиск заказа</h1>
			</div>
			<div>
				<input
					type='text'
					placeholder='Введите номер заказа'
					className='border border-gray-300 rounded p-2 mt-4 w-full max-w-md'
					value={orderId}
					onChange={handleInputChange}
				/>
			</div>
			<div className='mt-4'>
				<button
					className='bg-blue-500 text-white px-4 py-2 rounded'
					onClick={handleSearch}
					disabled={loading || !orderId}
				>
					{loading ? 'Поиск...' : 'Найти'}
				</button>
			</div>
			<div className='mt-4'>
				{error && <p className='text-red-500'>{error}</p>}
				{order ? (
					<div className='border p-4 rounded '>
						<h2 className='text-xl font-bold mb-2'>Информация о заказе</h2>
						<p>
							<b>Имя клиента:</b> {order.delivery?.name}
						</p>
						<p>
							<b>Сумма:</b> {order.payment?.amount} ₽
						</p>
						<p>
							<b>Товар:</b> {order.items?.name}
						</p>
						<p>
							<b>Доставка в:</b> {order.delivery?.city}
						</p>
					</div>
				) : (
					<p>Результаты поиска будут отображаться здесь.</p>
				)}
			</div>
		</>
	)
}

export default App
