import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';

const Footer = () => {
  return (
    <footer>
      <Container>
        <Row>
          <Col className="text-center py-3">
            Copyright &copy; 電子商務微服務平台 {new Date().getFullYear()}
          </Col>
        </Row>
      </Container>
    </footer>
  );
};

export default Footer;
