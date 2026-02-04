import React, { useEffect, useState } from 'react';
import {
  SafeAreaView,
  View,
  Text,
  TextInput,
  Button,
  FlatList,
  TouchableOpacity,
  StyleSheet,
  ActivityIndicator,
  ScrollView,
} from 'react-native';

// ВАЖНО: для эмулятора/реального девайса замените localhost
// Например, на IP вашей машины в локальной сети.
const AUTH_URL = 'http://localhost:8080';
const CATALOG_URL = 'http://localhost:8082';
const ORDER_URL = 'http://localhost:8085';

export default function App() {
  const [screen, setScreen] = useState('catalog'); // 'login' | 'catalog' | 'product' | 'cart' | 'orders'
  const [token, setToken] = useState(null);

  const [products, setProducts] = useState([]);
  const [loadingProducts, setLoadingProducts] = useState(false);
  const [productsError, setProductsError] = useState(null);

  const [selectedProduct, setSelectedProduct] = useState(null);

  const [cart, setCart] = useState([]); // [{ product, quantity }]

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [authError, setAuthError] = useState(null);

  const [orders, setOrders] = useState([]);
  const [ordersLoading, setOrdersLoading] = useState(false);
  const [ordersError, setOrdersError] = useState(null);

  const [infoMessage, setInfoMessage] = useState(null);

  useEffect(() => {
    loadProducts();
  }, []);

  async function loadProducts() {
    setLoadingProducts(true);
    setProductsError(null);
    try {
      const res = await fetch(`${CATALOG_URL}/products?status=active&limit=20&offset=0`);
      if (!res.ok) {
        throw new Error('Не удалось загрузить каталог');
      }
      const data = await res.json();
      setProducts(Array.isArray(data) ? data : []);
    } catch (e) {
      setProductsError(e.message || 'Ошибка загрузки каталога');
    } finally {
      setLoadingProducts(false);
    }
  }

  async function handleLogin() {
    setAuthError(null);
    setInfoMessage(null);
    try {
      if (!email || !password) {
        setAuthError('Введите email и пароль');
        return;
      }
      const res = await fetch(`${AUTH_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json();
      if (!res.ok) {
        setAuthError(data?.error || 'Ошибка входа');
        return;
      }
      if (!data.access_token) {
        setAuthError('Сервер не вернул access_token');
        return;
      }
      setToken(data.access_token);
      setScreen('catalog');
      setInfoMessage('Вы успешно вошли');
    } catch (e) {
      setAuthError(e.message || 'Ошибка сети при входе');
    }
  }

  function handleLogout() {
    setToken(null);
    setInfoMessage('Вы вышли из аккаунта');
  }

  function openProduct(product) {
    setSelectedProduct(product);
    setScreen('product');
  }

  function addToCart(product) {
    setCart((prev) => {
      const existing = prev.find((item) => item.product.id === product.id);
      if (existing) {
        return prev.map((item) =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        );
      }
      return [...prev, { product, quantity: 1 }];
    });
    setInfoMessage('Товар добавлен в корзину');
  }

  function removeFromCart(productId) {
    setCart((prev) => prev.filter((item) => item.product.id !== productId));
  }

  function changeQuantity(productId, delta) {
    setCart((prev) =>
      prev
        .map((item) =>
          item.product.id === productId
            ? { ...item, quantity: Math.max(1, item.quantity + delta) }
            : item,
        )
        .filter((item) => item.quantity > 0),
    );
  }

  const cartTotal = cart.reduce(
    (sum, item) => sum + (item.product.price || 0) * item.quantity,
    0,
  );

  async function placeOrder() {
    setInfoMessage(null);
    if (!token) {
      setScreen('login');
      setAuthError('Авторизуйтесь, чтобы оформить заказ');
      return;
    }
    if (cart.length === 0) {
      setInfoMessage('Корзина пуста');
      return;
    }
    try {
      const body = {
        items: cart.map((item) => ({
          product_id: item.product.id,
          name: item.product.name,
          price: item.product.price,
          quantity: item.quantity,
        })),
      };
      const res = await fetch(`${ORDER_URL}/orders`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });
      const data = await res.json();
      if (!res.ok) {
        setInfoMessage(data?.error || 'Не удалось создать заказ');
        return;
      }
      setCart([]);
      setInfoMessage(`Заказ создан (id: ${data.id || '—'})`);
      // можно сразу обновить список заказов
      await loadOrders();
      setScreen('orders');
    } catch (e) {
      setInfoMessage(e.message || 'Ошибка сети при создании заказа');
    }
  }

  async function loadOrders() {
    if (!token) {
      setScreen('login');
      setAuthError('Авторизуйтесь, чтобы просматривать заказы');
      return;
    }
    setOrdersLoading(true);
    setOrdersError(null);
    try {
      const res = await fetch(`${ORDER_URL}/orders`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      if (!res.ok) {
        setOrdersError(data?.error || 'Не удалось загрузить заказы');
        return;
      }
      setOrders(Array.isArray(data) ? data : []);
    } catch (e) {
      setOrdersError(e.message || 'Ошибка сети при загрузке заказов');
    } finally {
      setOrdersLoading(false);
    }
  }

  function renderLogin() {
    return (
      <View style={styles.block}>
        <Text style={styles.blockTitle}>Вход</Text>
        <TextInput
          style={styles.input}
          placeholder="Email"
          value={email}
          onChangeText={setEmail}
          autoCapitalize="none"
          keyboardType="email-address"
        />
        <TextInput
          style={styles.input}
          placeholder="Пароль"
          value={password}
          onChangeText={setPassword}
          secureTextEntry
        />
        {authError ? <Text style={styles.error}>{authError}</Text> : null}
        <Button title="Войти" onPress={handleLogin} />
        <Text style={styles.hint}>
          Сервер Auth: {AUTH_URL}
        </Text>
      </View>
    );
  }

  function renderCatalog() {
    if (loadingProducts) {
      return <ActivityIndicator style={styles.center} />;
    }
    if (productsError) {
      return (
        <View style={styles.block}>
          <Text style={styles.error}>{productsError}</Text>
          <Button title="Повторить" onPress={loadProducts} />
        </View>
      );
    }
    return (
      <View style={styles.block}>
        <Text style={styles.blockTitle}>Каталог</Text>
        <FlatList
          data={products}
          keyExtractor={(item) => String(item.id)}
          renderItem={({ item }) => (
            <TouchableOpacity
              style={styles.productCard}
              onPress={() => openProduct(item)}
            >
              <Text style={styles.productName}>{item.name}</Text>
              <Text style={styles.productPrice}>
                {item.price != null ? `${item.price} ₽` : '—'}
              </Text>
              <Button title="В корзину" onPress={() => addToCart(item)} />
            </TouchableOpacity>
          )}
          ListEmptyComponent={
            <Text style={styles.hint}>Товаров пока нет</Text>
          }
        />
      </View>
    );
  }

  function renderProduct() {
    if (!selectedProduct) {
      return (
        <View style={styles.block}>
          <Text>Товар не выбран</Text>
        </View>
      );
    }
    return (
      <ScrollView style={styles.block}>
        <Text style={styles.blockTitle}>{selectedProduct.name}</Text>
        {selectedProduct.description ? (
          <Text style={styles.productDescription}>
            {selectedProduct.description}
          </Text>
        ) : null}
        <Text style={styles.productPriceLarge}>
          {selectedProduct.price != null
            ? `${selectedProduct.price} ₽`
            : '—'}
        </Text>
        <Button title="В корзину" onPress={() => addToCart(selectedProduct)} />
      </ScrollView>
    );
  }

  function renderCart() {
    return (
      <View style={styles.block}>
        <Text style={styles.blockTitle}>Корзина</Text>
        <FlatList
          data={cart}
          keyExtractor={(item) => String(item.product.id)}
          renderItem={({ item }) => (
            <View style={styles.cartRow}>
              <View style={{ flex: 1 }}>
                <Text style={styles.productName}>{item.product.name}</Text>
                <Text style={styles.hint}>
                  {item.product.price != null ? `${item.product.price} ₽` : '—'}
                </Text>
              </View>
              <View style={styles.cartControls}>
                <Button
                  title="-"
                  onPress={() => changeQuantity(item.product.id, -1)}
                />
                <Text style={styles.cartQty}>{item.quantity}</Text>
                <Button
                  title="+"
                  onPress={() => changeQuantity(item.product.id, 1)}
                />
              </View>
              <Button
                title="Убрать"
                onPress={() => removeFromCart(item.product.id)}
              />
            </View>
          )}
          ListEmptyComponent={
            <Text style={styles.hint}>Корзина пуста</Text>
          }
        />
        <Text style={styles.total}>Итого: {cartTotal} ₽</Text>
        <Button title="Оформить заказ" onPress={placeOrder} />
      </View>
    );
  }

  function renderOrders() {
    return (
      <View style={styles.block}>
        <Text style={styles.blockTitle}>Мои заказы</Text>
        {ordersLoading ? <ActivityIndicator style={styles.center} /> : null}
        {ordersError ? <Text style={styles.error}>{ordersError}</Text> : null}
        <FlatList
          data={orders}
          keyExtractor={(item) => String(item.id)}
          renderItem={({ item }) => (
            <View style={styles.orderCard}>
              <Text style={styles.productName}>Заказ {item.id}</Text>
              <Text style={styles.hint}>Статус: {item.status}</Text>
              <Text style={styles.hint}>Сумма: {item.total} ₽</Text>
            </View>
          )}
          ListEmptyComponent={
            !ordersLoading && (
              <Text style={styles.hint}>Заказов пока нет</Text>
            )
          }
        />
      </View>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Marketplace Mobile</Text>
      </View>

      <View style={styles.nav}>
        <Button title="Каталог" onPress={() => setScreen('catalog')} />
        <Button title="Корзина" onPress={() => setScreen('cart')} />
        <Button
          title="Заказы"
          onPress={() => {
            if (!token) {
              setScreen('login');
            } else {
              loadOrders();
              setScreen('orders');
            }
          }}
        />
        <Button
          title={token ? 'Выход' : 'Войти'}
          onPress={() => {
            if (token) {
              handleLogout();
            } else {
              setScreen('login');
            }
          }}
        />
      </View>

      {infoMessage ? <Text style={styles.info}>{infoMessage}</Text> : null}

      <View style={styles.content}>
        {screen === 'login' && renderLogin()}
        {screen === 'catalog' && renderCatalog()}
        {screen === 'product' && renderProduct()}
        {screen === 'cart' && renderCart()}
        {screen === 'orders' && renderOrders()}
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    padding: 16,
    backgroundColor: '#2563eb',
  },
  title: {
    color: '#fff',
    fontSize: 20,
    fontWeight: '600',
  },
  nav: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingVertical: 8,
    backgroundColor: '#e5e7eb',
  },
  content: {
    flex: 1,
    padding: 12,
  },
  block: {
    backgroundColor: '#fff',
    borderRadius: 8,
    padding: 12,
    flex: 1,
  },
  blockTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 8,
  },
  input: {
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 6,
    paddingHorizontal: 10,
    paddingVertical: 8,
    marginBottom: 8,
    backgroundColor: '#f9fafb',
  },
  error: {
    color: '#b91c1c',
    marginBottom: 8,
  },
  info: {
    color: '#065f46',
    paddingHorizontal: 12,
    paddingVertical: 4,
  },
  hint: {
    color: '#6b7280',
    marginTop: 4,
  },
  productCard: {
    backgroundColor: '#fff',
    padding: 10,
    borderRadius: 8,
    marginBottom: 8,
  },
  productName: {
    fontSize: 16,
    fontWeight: '500',
  },
  productPrice: {
    marginTop: 4,
    marginBottom: 4,
    fontWeight: '600',
  },
  productPriceLarge: {
    fontSize: 22,
    fontWeight: '700',
    marginVertical: 12,
  },
  productDescription: {
    marginTop: 8,
    marginBottom: 8,
    color: '#374151',
  },
  cartRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
    backgroundColor: '#fff',
    padding: 8,
    borderRadius: 8,
  },
  cartControls: {
    flexDirection: 'row',
    alignItems: 'center',
    marginHorizontal: 8,
  },
  cartQty: {
    marginHorizontal: 8,
    fontWeight: '600',
  },
  total: {
    marginTop: 8,
    fontSize: 16,
    fontWeight: '700',
  },
  orderCard: {
    backgroundColor: '#fff',
    padding: 10,
    borderRadius: 8,
    marginBottom: 8,
  },
  center: {
    marginTop: 16,
  },
});

