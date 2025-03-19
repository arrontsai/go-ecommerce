import React, { useEffect, useState } from 'react';
import { Link, useNavigate, useParams, useLocation } from 'react-router-dom';
import { Row, Col, ListGroup, Image, Form, Button, Card } from 'react-bootstrap';
import axios from 'axios';

const CartPage = () => {
  const { id } = useParams();
  const location = useLocation();
  const navigate = useNavigate();

  const qty = location.search ? Number(location.search.split('=')[1]) : 1;
  const [cartItems, setCartItems] = useState([]);

  useEffect(() => {
    // 如果有產品 ID，則添加到購物車
    const addToCart = async () => {
      if (id) {
        try {
          const { data } = await axios.get(`http://localhost:8082/api/v1/products/${id}`);
          const item = {
            product: data.id,
            name: data.name,
            image: data.imageUrl,
            price: data.price,
            countInStock: data.countInStock,
            qty,
          };
          
          // 檢查購物車中是否已有該產品
          const existItem = cartItems.find((x) => x.product === item.product);
          
          if (existItem) {
            setCartItems(
              cartItems.map((x) =>
                x.product === existItem.product ? item : x
              )
            );
          } else {
            setCartItems([...cartItems, item]);
          }
          
          // 保存到 localStorage
          localStorage.setItem('cartItems', JSON.stringify([...cartItems, item]));
        } catch (error) {
          console.error('Error fetching product:', error);
          
          // 開發測試用模擬數據
          const item = {
            product: id,
            name: '模擬商品',
            image: 'https://via.placeholder.com/100',
            price: 199,
            countInStock: 5,
            qty,
          };
          
          const existItem = cartItems.find((x) => x.product === item.product);
          
          if (existItem) {
            setCartItems(
              cartItems.map((x) =>
                x.product === existItem.product ? item : x
              )
            );
          } else {
            setCartItems([...cartItems, item]);
          }
          
          localStorage.setItem('cartItems', JSON.stringify([...cartItems, item]));
        }
      } else {
        // 從 localStorage 讀取購物車數據
        const storedCartItems = localStorage.getItem('cartItems');
        if (storedCartItems) {
          setCartItems(JSON.parse(storedCartItems));
        }
      }
    };

    addToCart();
  }, [id, qty]);

  const removeFromCartHandler = (id) => {
    const updatedCartItems = cartItems.filter((item) => item.product !== id);
    setCartItems(updatedCartItems);
    localStorage.setItem('cartItems', JSON.stringify(updatedCartItems));
  };

  const checkoutHandler = () => {
    // 檢查是否已登入
    const userInfo = localStorage.getItem('userInfo');
    if (!userInfo) {
      navigate('/login?redirect=shipping');
    } else {
      navigate('/shipping');
    }
  };

  return (
    <Row>
      <Col md={8}>
        <h1>購物車</h1>
        {cartItems.length === 0 ? (
          <div className="alert alert-info">
            您的購物車是空的 <Link to="/">返回</Link>
          </div>
        ) : (
          <ListGroup variant="flush">
            {cartItems.map((item) => (
              <ListGroup.Item key={item.product}>
                <Row>
                  <Col md={2}>
                    <Image src={item.image} alt={item.name} fluid rounded />
                  </Col>
                  <Col md={3}>
                    <Link to={`/product/${item.product}`}>{item.name}</Link>
                  </Col>
                  <Col md={2}>NT${item.price}</Col>
                  <Col md={2}>
                    <Form.Control
                      as="select"
                      value={item.qty}
                      onChange={(e) => {
                        const updatedCartItems = cartItems.map((x) =>
                          x.product === item.product
                            ? { ...x, qty: Number(e.target.value) }
                            : x
                        );
                        setCartItems(updatedCartItems);
                        localStorage.setItem('cartItems', JSON.stringify(updatedCartItems));
                      }}
                    >
                      {[...Array(item.countInStock).keys()].map((x) => (
                        <option key={x + 1} value={x + 1}>
                          {x + 1}
                        </option>
                      ))}
                    </Form.Control>
                  </Col>
                  <Col md={2}>
                    <Button
                      type="button"
                      variant="light"
                      onClick={() => removeFromCartHandler(item.product)}
                    >
                      <i className="fas fa-trash"></i>
                    </Button>
                  </Col>
                </Row>
              </ListGroup.Item>
            ))}
          </ListGroup>
        )}
      </Col>
      <Col md={4}>
        <Card>
          <ListGroup variant="flush">
            <ListGroup.Item>
              <h2>
                小計 ({cartItems.reduce((acc, item) => acc + item.qty, 0)}) 件商品
              </h2>
              NT$
              {cartItems
                .reduce((acc, item) => acc + item.qty * item.price, 0)
                .toFixed(2)}
            </ListGroup.Item>
            <ListGroup.Item>
              <Button
                type="button"
                className="btn-block"
                disabled={cartItems.length === 0}
                onClick={checkoutHandler}
              >
                前往結帳
              </Button>
            </ListGroup.Item>
          </ListGroup>
        </Card>
      </Col>
    </Row>
  );
};

export default CartPage;
