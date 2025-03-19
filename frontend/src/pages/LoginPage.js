import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Form, Button, Row, Col } from 'react-bootstrap';
import axios from 'axios';

const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  
  const navigate = useNavigate();

  useEffect(() => {
    // 如果已經登入，導向主頁
    const userInfo = localStorage.getItem('userInfo');
    if (userInfo) {
      navigate('/');
    }
  }, [navigate]);

  const submitHandler = async (e) => {
    e.preventDefault();
    try {
      const config = {
        headers: {
          'Content-Type': 'application/json',
        },
      };

      const { data } = await axios.post(
        'http://localhost:8081/api/v1/auth/login',
        { email, password },
        config
      );

      localStorage.setItem('userInfo', JSON.stringify(data));
      navigate('/');
    } catch (error) {
      setError(
        error.response && error.response.data.message
          ? error.response.data.message
          : '登入失敗，請確認您的電子郵件和密碼'
      );
    }
  };

  return (
    <>
      <h1>登入</h1>
      {error && <div className="alert alert-danger">{error}</div>}
      <Form onSubmit={submitHandler}>
        <Form.Group controlId="email" className="my-3">
          <Form.Label>電子郵件</Form.Label>
          <Form.Control
            type="email"
            placeholder="請輸入電子郵件"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </Form.Group>

        <Form.Group controlId="password" className="my-3">
          <Form.Label>密碼</Form.Label>
          <Form.Control
            type="password"
            placeholder="請輸入密碼"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </Form.Group>

        <Button type="submit" variant="primary" className="my-3">
          登入
        </Button>
      </Form>

      <Row className="py-3">
        <Col>
          還沒有帳號？{' '}
          <Link to="/register">註冊</Link>
        </Col>
      </Row>
    </>
  );
};

export default LoginPage;
