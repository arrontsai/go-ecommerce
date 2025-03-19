import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Form, Button, Row, Col } from 'react-bootstrap';
import axios from 'axios';

const RegisterPage = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
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
    
    if (password !== confirmPassword) {
      setError('密碼不匹配');
      return;
    }
    
    try {
      const config = {
        headers: {
          'Content-Type': 'application/json',
        },
      };

      const { data } = await axios.post(
        'http://localhost:8081/api/v1/auth/register',
        { name, email, password },
        config
      );

      localStorage.setItem('userInfo', JSON.stringify(data));
      navigate('/');
    } catch (error) {
      setError(
        error.response && error.response.data.message
          ? error.response.data.message
          : '註冊失敗，請稍後再試'
      );
    }
  };

  return (
    <>
      <h1>註冊</h1>
      {error && <div className="alert alert-danger">{error}</div>}
      <Form onSubmit={submitHandler}>
        <Form.Group controlId="name" className="my-3">
          <Form.Label>姓名</Form.Label>
          <Form.Control
            type="text"
            placeholder="請輸入姓名"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </Form.Group>

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

        <Form.Group controlId="confirmPassword" className="my-3">
          <Form.Label>確認密碼</Form.Label>
          <Form.Control
            type="password"
            placeholder="請再次輸入密碼"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
        </Form.Group>

        <Button type="submit" variant="primary" className="my-3">
          註冊
        </Button>
      </Form>

      <Row className="py-3">
        <Col>
          已經有帳號？{' '}
          <Link to="/login">登入</Link>
        </Col>
      </Row>
    </>
  );
};

export default RegisterPage;
