// OmniRoute Shop Portal - Homepage
// B2C Consumer E-commerce Homepage

'use client';

import { useState } from 'react';
import { Input, Button, Card, Row, Col, Typography, Tag, Badge, Space, Carousel, Avatar, Rate } from 'antd';
import {
    SearchOutlined,
    ShoppingCartOutlined,
    HeartOutlined,
    RightOutlined,
    FireOutlined,
    ThunderboltOutlined,
    TruckOutlined
} from '@ant-design/icons';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

// Types
interface Product {
    id: string;
    name: string;
    slug: string;
    price: number;
    comparePrice?: number;
    image: string;
    category: string;
    rating: number;
    reviews: number;
    isNew?: boolean;
    isBestSeller?: boolean;
}

interface Category {
    id: string;
    name: string;
    slug: string;
    icon: string;
    productCount: number;
}

// Mock Data
const featuredProducts: Product[] = [
    { id: '1', name: 'Peak Evaporated Milk 400g', slug: 'peak-milk-400g', price: 2800, comparePrice: 3200, image: '/products/peak-milk.jpg', category: 'Dairy', rating: 4.8, reviews: 234, isBestSeller: true },
    { id: '2', name: 'Golden Penny Rice 50kg', slug: 'golden-penny-rice-50kg', price: 68000, image: '/products/golden-penny-rice.jpg', category: 'Grains', rating: 4.9, reviews: 567, isBestSeller: true },
    { id: '3', name: 'Indomie Chicken Flavour (Carton)', slug: 'indomie-chicken-carton', price: 5500, comparePrice: 6000, image: '/products/indomie.jpg', category: 'Noodles', rating: 4.7, reviews: 890, isNew: true },
    { id: '4', name: 'Kings Vegetable Oil 5L', slug: 'kings-oil-5l', price: 9800, image: '/products/kings-oil.jpg', category: 'Cooking Oil', rating: 4.6, reviews: 123 },
    { id: '5', name: 'Dangote Sugar 50kg', slug: 'dangote-sugar-50kg', price: 48000, image: '/products/dangote-sugar.jpg', category: 'Sugar', rating: 4.8, reviews: 345 },
    { id: '6', name: 'Milo Tin 500g', slug: 'milo-tin-500g', price: 3200, comparePrice: 3500, image: '/products/milo.jpg', category: 'Beverages', rating: 4.9, reviews: 678 },
];

const categories: Category[] = [
    { id: '1', name: 'Beverages', slug: 'beverages', icon: '🥤', productCount: 245 },
    { id: '2', name: 'Food Items', slug: 'food-items', icon: '🍚', productCount: 567 },
    { id: '3', name: 'Personal Care', slug: 'personal-care', icon: '🧴', productCount: 189 },
    { id: '4', name: 'Household', slug: 'household', icon: '🏠', productCount: 234 },
    { id: '5', name: 'Baby Products', slug: 'baby-products', icon: '👶', productCount: 145 },
    { id: '6', name: 'Electronics', slug: 'electronics', icon: '📱', productCount: 89 },
];

const bannerSlides = [
    { id: '1', title: 'New Year Mega Sale!', subtitle: 'Up to 40% off on select items', cta: 'Shop Now', bg: 'linear-gradient(135deg, #1a365d 0%, #2d3748 100%)' },
    { id: '2', title: 'Free Delivery', subtitle: 'On orders above ₦50,000', cta: 'Learn More', bg: 'linear-gradient(135deg, #2f855a 0%, #276749 100%)' },
    { id: '3', title: 'Bulk Discounts', subtitle: 'Save more when you buy more', cta: 'View Deals', bg: 'linear-gradient(135deg, #d69e2e 0%, #b7791f 100%)' },
];

// Formatters
const formatNaira = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
    }).format(amount);
};

// Components
const ProductCard = ({ product }: { product: Product }) => (
    <Card
        hoverable
        cover={
            <div style={{
                height: 200,
                background: '#f5f5f5',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'relative'
            }}>
                <Text type="secondary">Product Image</Text>
                {product.isNew && (
                    <Tag color="blue" style={{ position: 'absolute', top: 8, left: 8 }}>NEW</Tag>
                )}
                {product.isBestSeller && (
                    <Tag color="red" icon={<FireOutlined />} style={{ position: 'absolute', top: 8, left: 8 }}>
                        Best Seller
                    </Tag>
                )}
                <Button
                    shape="circle"
                    icon={<HeartOutlined />}
                    style={{ position: 'absolute', top: 8, right: 8 }}
                />
            </div>
        }
        bodyStyle={{ padding: 12 }}
    >
        <Text type="secondary" style={{ fontSize: 12 }}>{product.category}</Text>
        <Paragraph ellipsis={{ rows: 2 }} style={{ marginBottom: 4, height: 44 }}>
            {product.name}
        </Paragraph>
        <Space size={4}>
            <Rate disabled defaultValue={product.rating} style={{ fontSize: 12 }} />
            <Text type="secondary" style={{ fontSize: 12 }}>({product.reviews})</Text>
        </Space>
        <div style={{ marginTop: 8 }}>
            <Text strong style={{ fontSize: 18, color: '#1890ff' }}>{formatNaira(product.price)}</Text>
            {product.comparePrice && (
                <Text delete type="secondary" style={{ marginLeft: 8, fontSize: 14 }}>
                    {formatNaira(product.comparePrice)}
                </Text>
            )}
        </div>
        <Button type="primary" icon={<ShoppingCartOutlined />} block style={{ marginTop: 12 }}>
            Add to Cart
        </Button>
    </Card>
);

// Main Component
export default function ShopHomepage() {
    const [cartCount, setCartCount] = useState(3);

    return (
        <div>
            {/* Header */}
            <header style={{
                background: '#1a365d',
                padding: '16px 24px',
                position: 'sticky',
                top: 0,
                zIndex: 100
            }}>
                <Row align="middle" gutter={24}>
                    <Col>
                        <Title level={3} style={{ color: 'white', margin: 0 }}>OmniRoute</Title>
                    </Col>
                    <Col flex="auto">
                        <Search
                            placeholder="Search products..."
                            size="large"
                            enterButton={<SearchOutlined />}
                            style={{ maxWidth: 600 }}
                        />
                    </Col>
                    <Col>
                        <Space size="large">
                            <Button type="text" style={{ color: 'white' }}>Categories</Button>
                            <Button type="text" style={{ color: 'white' }}>Deals</Button>
                            <Badge count={cartCount}>
                                <Button type="primary" icon={<ShoppingCartOutlined />}>
                                    Cart
                                </Button>
                            </Badge>
                            <Avatar icon={<Avatar>JO</Avatar>} />
                        </Space>
                    </Col>
                </Row>
            </header>

            {/* Hero Banner */}
            <Carousel autoplay effect="fade">
                {bannerSlides.map((slide) => (
                    <div key={slide.id}>
                        <div style={{
                            background: slide.bg,
                            padding: '60px 24px',
                            textAlign: 'center'
                        }}>
                            <Title style={{ color: 'white', marginBottom: 8 }}>{slide.title}</Title>
                            <Paragraph style={{ color: 'rgba(255,255,255,0.8)', fontSize: 18, marginBottom: 24 }}>
                                {slide.subtitle}
                            </Paragraph>
                            <Button type="primary" size="large" ghost>{slide.cta} <RightOutlined /></Button>
                        </div>
                    </div>
                ))}
            </Carousel>

            {/* Value Props */}
            <div style={{ background: '#f5f5f5', padding: '16px 24px' }}>
                <Row gutter={24} justify="center">
                    <Col>
                        <Space>
                            <TruckOutlined style={{ fontSize: 24, color: '#1890ff' }} />
                            <Text strong>Free Delivery on ₦50K+</Text>
                        </Space>
                    </Col>
                    <Col>
                        <Space>
                            <ThunderboltOutlined style={{ fontSize: 24, color: '#52c41a' }} />
                            <Text strong>Same-Day Delivery</Text>
                        </Space>
                    </Col>
                    <Col>
                        <Space>
                            <ShoppingCartOutlined style={{ fontSize: 24, color: '#fa8c16' }} />
                            <Text strong>Bulk Discounts Available</Text>
                        </Space>
                    </Col>
                </Row>
            </div>

            <div style={{ padding: '32px 24px', maxWidth: 1400, margin: '0 auto' }}>
                {/* Categories */}
                <div style={{ marginBottom: 48 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
                        <Title level={3} style={{ margin: 0 }}>Shop by Category</Title>
                        <Button type="link">View All <RightOutlined /></Button>
                    </div>
                    <Row gutter={[16, 16]}>
                        {categories.map((cat) => (
                            <Col xs={12} sm={8} md={4} key={cat.id}>
                                <Card hoverable style={{ textAlign: 'center' }}>
                                    <div style={{ fontSize: 32, marginBottom: 8 }}>{cat.icon}</div>
                                    <Text strong>{cat.name}</Text>
                                    <br />
                                    <Text type="secondary" style={{ fontSize: 12 }}>{cat.productCount} items</Text>
                                </Card>
                            </Col>
                        ))}
                    </Row>
                </div>

                {/* Featured Products */}
                <div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
                        <Title level={3} style={{ margin: 0 }}>
                            <FireOutlined style={{ color: '#ff4d4f', marginRight: 8 }} />
                            Featured Products
                        </Title>
                        <Button type="link">View All <RightOutlined /></Button>
                    </div>
                    <Row gutter={[16, 16]}>
                        {featuredProducts.map((product) => (
                            <Col xs={24} sm={12} md={8} lg={4} key={product.id}>
                                <ProductCard product={product} />
                            </Col>
                        ))}
                    </Row>
                </div>
            </div>

            {/* Footer */}
            <footer style={{ background: '#1a365d', color: 'white', padding: '48px 24px 24px' }}>
                <Row gutter={48} style={{ maxWidth: 1400, margin: '0 auto' }}>
                    <Col xs={24} md={8}>
                        <Title level={4} style={{ color: 'white' }}>OmniRoute</Title>
                        <Paragraph style={{ color: 'rgba(255,255,255,0.7)' }}>
                            Nigeria's leading B2B/B2C FMCG commerce platform. Connecting manufacturers, distributors, and retailers.
                        </Paragraph>
                    </Col>
                    <Col xs={12} md={4}>
                        <Title level={5} style={{ color: 'white' }}>Shop</Title>
                        <Space direction="vertical">
                            <a style={{ color: 'rgba(255,255,255,0.7)' }}>All Products</a>
                            <a style={{ color: 'rgba(255,255,255,0.7)' }}>Categories</a>
                            <a style={{ color: 'rgba(255,255,255,0.7)' }}>Deals</a>
                        </Space>
                    </Col>
                    <Col xs={12} md={4}>
                        <Title level={5} style={{ color: 'white' }}>Account</Title>
                        <Space direction="vertical">
                            <a style={{ color: 'rgba(255,255,255,0.7)' }}>My Orders</a>
                            <a style={{ color: 'rgba(255,255,255,0.7)' }}>Wishlist</a>
                            <a style={{ color: 'rgba(255,255,255,0.7)' }}>Profile</a>
                        </Space>
                    </Col>
                    <Col xs={24} md={8}>
                        <Title level={5} style={{ color: 'white' }}>Contact Us</Title>
                        <Paragraph style={{ color: 'rgba(255,255,255,0.7)' }}>
                            📞 0800-OMNIROUTE<br />
                            📧 support@omniroute.io<br />
                            💬 WhatsApp: +234 800 123 4567
                        </Paragraph>
                    </Col>
                </Row>
                <div style={{ textAlign: 'center', marginTop: 48, paddingTop: 24, borderTop: '1px solid rgba(255,255,255,0.1)' }}>
                    <Text style={{ color: 'rgba(255,255,255,0.5)' }}>
                        © 2026 OmniRoute Commerce Platform. All rights reserved.
                    </Text>
                </div>
            </footer>
        </div>
    );
}
