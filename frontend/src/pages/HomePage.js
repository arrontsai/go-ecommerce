import React, { useState, useEffect } from 'react';
import { Row, Col, Card } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import axios from 'axios';

const HomePage = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const { data } = await axios.get('http://localhost:8082/api/v1/products');
        setProducts(data.products || []);
        setLoading(false);
      } catch (error) {
        setError(error.message);
        setLoading(false);
        // 開發測試用模擬數據
        setProducts([
          {
            id: '1',
            name: '商品一',
            description: '這是商品一的描述',
            price: 100,
            imageUrl: 'https://via.placeholder.com/150',
          },
          {
            id: '2',
            name: '商品二',
            description: '這是商品二的描述',
            price: 200,
            imageUrl: 'https://via.placeholder.com/150',
          },
          {
            id: '3',
            name: '商品三',
            description: '這是商品三的描述',
            price: 300,
            imageUrl: 'https://via.placeholder.com/150',
          },
        ]);
      }
    };

    fetchProducts();
  }, []);

  return (
    <>
      <h1>最新商品</h1>
      {loading ? (
        <p>載入中...</p>
      ) : error ? (
        <p>發生錯誤：{error}</p>
      ) : (
        <Row>
          {products.map((product) => (
            <Col key={product.id} sm={12} md={6} lg={4} xl={3}>
              <Card className="my-3 p-3 rounded product-card">
                <Link to={`/product/${product.id}`}>
                  <Card.Img src={product.imageUrl} variant="top" />
                </Link>

                <Card.Body>
                  <Link to={`/product/${product.id}`}>
                    <Card.Title as="div" className="product-title">
                      <strong>{product.name}</strong>
                    </Card.Title>
                  </Link>

                  <Card.Text as="h3">NT${product.price}</Card.Text>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </>
  );
};

export default HomePage;
