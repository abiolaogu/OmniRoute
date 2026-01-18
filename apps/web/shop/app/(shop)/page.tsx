// OmniRoute Shop Portal - Innovative B2C E-commerce Experience
// AI-Powered Shopping, Voice Search, AR Preview, Social Commerce

'use client';

import React, { useState, useEffect, useRef } from 'react';
import {
    Input, Button, Card, Row, Col, Typography, Tag, Badge, Space,
    Carousel, Avatar, Rate, Drawer, List, Divider, Tabs, Segmented,
    Tooltip, Modal, Progress, notification, Affix, FloatButton, Popover,
    Skeleton, Empty, Alert
} from 'antd';
import {
    SearchOutlined, ShoppingCartOutlined, HeartOutlined, HeartFilled,
    RightOutlined, FireOutlined, ThunderboltOutlined, TruckOutlined,
    GiftOutlined, StarOutlined, StarFilled, UserOutlined, HomeOutlined,
    AppstoreOutlined, HistoryOutlined, BellOutlined, CustomerServiceOutlined,
    AudioOutlined, CameraOutlined, ScanOutlined, ShareAltOutlined,
    EnvironmentOutlined, ClockCircleOutlined, SafetyCertificateOutlined,
    WalletOutlined, CreditCardOutlined, BankOutlined, TeamOutlined,
    TrophyOutlined, RobotOutlined, EyeOutlined, PlusOutlined, MinusOutlined,
    DeleteOutlined, WhatsAppOutlined, SendOutlined, PercentageOutlined,
    LikeOutlined, MessageOutlined, ShopOutlined
} from '@ant-design/icons';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

// =============================================================================
// TYPES
// =============================================================================

interface Product {
    id: string;
    name: string;
    slug: string;
    price: number;
    comparePrice?: number;
    image: string;
    category: string;
    brand: string;
    rating: number;
    reviews: number;
    sold: number;
    stock: number;
    isNew?: boolean;
    isBestSeller?: boolean;
    isFlashSale?: boolean;
    flashSaleEnds?: string;
    discount?: number;
    variants?: { name: string; options: string[] }[];
    description?: string;
}

interface Category {
    id: string;
    name: string;
    slug: string;
    icon: string;
    productCount: number;
    image?: string;
}

interface CartItem {
    product: Product;
    quantity: number;
    variant?: string;
}

interface Store {
    id: string;
    name: string;
    rating: number;
    products: number;
    verified: boolean;
    image?: string;
}

// =============================================================================
// MOCK DATA
// =============================================================================

const featuredProducts: Product[] = [
    { id: '1', name: 'Peak Evaporated Milk 400g', slug: 'peak-milk-400g', price: 2800, comparePrice: 3200, image: '/products/peak-milk.jpg', category: 'Dairy', brand: 'Peak', rating: 4.8, reviews: 1234, sold: 5678, stock: 450, isBestSeller: true, discount: 12 },
    { id: '2', name: 'Golden Penny Rice 50kg (Premium)', slug: 'golden-penny-rice-50kg', price: 68000, comparePrice: 75000, image: '/products/golden-penny-rice.jpg', category: 'Grains', brand: 'Golden Penny', rating: 4.9, reviews: 567, sold: 2345, stock: 120, isBestSeller: true },
    { id: '3', name: 'Indomie Chicken Super Pack (40)', slug: 'indomie-chicken-40', price: 5500, image: '/products/indomie.jpg', category: 'Noodles', brand: 'Indomie', rating: 4.7, reviews: 890, sold: 12340, stock: 89, isNew: true },
    { id: '4', name: 'Kings Vegetable Oil 5L', slug: 'kings-oil-5l', price: 9800, comparePrice: 11000, image: '/products/kings-oil.jpg', category: 'Cooking Oil', brand: 'Kings', rating: 4.6, reviews: 345, sold: 890, stock: 234 },
    { id: '5', name: 'Dangote Sugar 50kg', slug: 'dangote-sugar-50kg', price: 48000, image: '/products/dangote-sugar.jpg', category: 'Sugar', brand: 'Dangote', rating: 4.8, reviews: 678, sold: 3456, stock: 320 },
    { id: '6', name: 'Milo Active-Go Tin 500g', slug: 'milo-tin-500g', price: 3200, comparePrice: 3800, image: '/products/milo.jpg', category: 'Beverages', brand: 'Nestle', rating: 4.9, reviews: 2345, sold: 8901, stock: 567, isFlashSale: true, flashSaleEnds: '2026-01-18T23:59:59', discount: 16 },
    { id: '7', name: 'Cowbell Milk Powder 400g', slug: 'cowbell-400g', price: 2500, image: '/products/cowbell.jpg', category: 'Dairy', brand: 'Cowbell', rating: 4.5, reviews: 432, sold: 2100, stock: 340 },
    { id: '8', name: 'Power Oil 5L', slug: 'power-oil-5l', price: 8500, comparePrice: 9500, image: '/products/power-oil.jpg', category: 'Cooking Oil', brand: 'Power', rating: 4.4, reviews: 567, sold: 1890, stock: 178, discount: 11 },
];

const categories: Category[] = [
    { id: '1', name: 'Beverages', slug: 'beverages', icon: '🥤', productCount: 245 },
    { id: '2', name: 'Food Items', slug: 'food-items', icon: '🍚', productCount: 567 },
    { id: '3', name: 'Personal Care', slug: 'personal-care', icon: '🧴', productCount: 189 },
    { id: '4', name: 'Household', slug: 'household', icon: '🏠', productCount: 234 },
    { id: '5', name: 'Baby Products', slug: 'baby-products', icon: '👶', productCount: 145 },
    { id: '6', name: 'Snacks', slug: 'snacks', icon: '🍿', productCount: 312 },
    { id: '7', name: 'Frozen Foods', slug: 'frozen-foods', icon: '🧊', productCount: 89 },
    { id: '8', name: 'Fresh Produce', slug: 'fresh-produce', icon: '🥬', productCount: 156 },
];

const popularStores: Store[] = [
    { id: '1', name: 'Lagos Fresh Mart', rating: 4.9, products: 234, verified: true },
    { id: '2', name: 'Kano Wholesale Hub', rating: 4.7, products: 567, verified: true },
    { id: '3', name: 'Abuja Premium Store', rating: 4.8, products: 189, verified: true },
];

// =============================================================================
// UTILITIES
// =============================================================================

const formatNaira = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
        style: 'currency',
        currency: 'NGN',
        minimumFractionDigits: 0,
    }).format(amount);
};

// =============================================================================
// COMPONENTS
// =============================================================================

const FlashSaleTimer = ({ endTime }: { endTime: string }) => {
    const [timeLeft, setTimeLeft] = useState({ hours: 0, minutes: 0, seconds: 0 });

    useEffect(() => {
        const timer = setInterval(() => {
            const diff = new Date(endTime).getTime() - Date.now();
            if (diff > 0) {
                setTimeLeft({
                    hours: Math.floor(diff / (1000 * 60 * 60)),
                    minutes: Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60)),
                    seconds: Math.floor((diff % (1000 * 60)) / 1000),
                });
            }
        }, 1000);
        return () => clearInterval(timer);
    }, [endTime]);

    return (
        <Space>
            <Tag color="red">{String(timeLeft.hours).padStart(2, '0')}</Tag>:
            <Tag color="red">{String(timeLeft.minutes).padStart(2, '0')}</Tag>:
            <Tag color="red">{String(timeLeft.seconds).padStart(2, '0')}</Tag>
        </Space>
    );
};

const ProductCard = ({ product, onAddToCart, onAddToWishlist }: {
    product: Product;
    onAddToCart: (p: Product) => void;
    onAddToWishlist: (p: Product) => void;
}) => {
    const [isWishlisted, setIsWishlisted] = useState(false);
    const [isHovered, setIsHovered] = useState(false);

    return (
        <Card
            hoverable
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
            cover={
                <div style={{
                    height: 180,
                    background: 'linear-gradient(135deg, #f5f5f5 0%, #e8e8e8 100%)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    position: 'relative',
                    overflow: 'hidden',
                    transition: 'all 0.3s ease'
                }}>
                    <Text type="secondary" style={{ fontSize: 48 }}>📦</Text>

                    {/* Badges */}
                    <div style={{ position: 'absolute', top: 8, left: 8, display: 'flex', flexDirection: 'column', gap: 4 }}>
                        {product.isFlashSale && (
                            <Tag color="red" icon={<ThunderboltOutlined />}>FLASH SALE</Tag>
                        )}
                        {product.isNew && <Tag color="blue">NEW</Tag>}
                        {product.isBestSeller && <Tag color="gold" icon={<FireOutlined />}>Best Seller</Tag>}
                        {product.discount && !product.isFlashSale && (
                            <Tag color="green">-{product.discount}%</Tag>
                        )}
                    </div>

                    {/* Quick Actions */}
                    <div style={{
                        position: 'absolute',
                        top: 8,
                        right: 8,
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 8,
                        opacity: isHovered ? 1 : 0,
                        transform: isHovered ? 'translateX(0)' : 'translateX(10px)',
                        transition: 'all 0.3s ease'
                    }}>
                        <Tooltip title={isWishlisted ? 'Remove from Wishlist' : 'Add to Wishlist'}>
                            <Button
                                shape="circle"
                                icon={isWishlisted ? <HeartFilled style={{ color: '#ff4d4f' }} /> : <HeartOutlined />}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    setIsWishlisted(!isWishlisted);
                                    onAddToWishlist(product);
                                }}
                            />
                        </Tooltip>
                        <Tooltip title="Quick View">
                            <Button shape="circle" icon={<EyeOutlined />} />
                        </Tooltip>
                        <Tooltip title="Share">
                            <Button shape="circle" icon={<ShareAltOutlined />} />
                        </Tooltip>
                    </div>

                    {/* Stock Warning */}
                    {product.stock < 50 && (
                        <div style={{ position: 'absolute', bottom: 0, left: 0, right: 0, background: 'rgba(255,77,79,0.9)', padding: '4px 8px' }}>
                            <Text style={{ color: 'white', fontSize: 11 }}>
                                🔥 Only {product.stock} left! Order now
                            </Text>
                        </div>
                    )}
                </div>
            }
            bodyStyle={{ padding: 12 }}
        >
            <div style={{ display: 'flex', alignItems: 'center', gap: 4, marginBottom: 4 }}>
                <Text type="secondary" style={{ fontSize: 11 }}>{product.brand}</Text>
                <Divider type="vertical" style={{ margin: 0 }} />
                <Text type="secondary" style={{ fontSize: 11 }}>{product.category}</Text>
            </div>

            <Paragraph ellipsis={{ rows: 2 }} style={{ marginBottom: 4, height: 44, lineHeight: 1.4 }}>
                <Text strong>{product.name}</Text>
            </Paragraph>

            <Space size={4} style={{ marginBottom: 8 }}>
                <Rate disabled defaultValue={product.rating} style={{ fontSize: 12 }} />
                <Text type="secondary" style={{ fontSize: 11 }}>({product.reviews.toLocaleString()})</Text>
                <Divider type="vertical" style={{ margin: 0 }} />
                <Text type="secondary" style={{ fontSize: 11 }}>{product.sold.toLocaleString()} sold</Text>
            </Space>

            <div style={{ marginBottom: 12 }}>
                <Text strong style={{ fontSize: 20, color: '#ff4d4f' }}>{formatNaira(product.price)}</Text>
                {product.comparePrice && (
                    <Text delete type="secondary" style={{ marginLeft: 8, fontSize: 13 }}>
                        {formatNaira(product.comparePrice)}
                    </Text>
                )}
            </div>

            {product.isFlashSale && product.flashSaleEnds && (
                <div style={{ marginBottom: 8 }}>
                    <Text type="secondary" style={{ fontSize: 11 }}>Ends in: </Text>
                    <FlashSaleTimer endTime={product.flashSaleEnds} />
                </div>
            )}

            <Button
                type="primary"
                icon={<ShoppingCartOutlined />}
                block
                onClick={(e) => {
                    e.stopPropagation();
                    onAddToCart(product);
                    notification.success({
                        message: 'Added to Cart',
                        description: `${product.name} has been added to your cart.`,
                        placement: 'bottomRight',
                        duration: 2,
                    });
                }}
            >
                Add to Cart
            </Button>
        </Card>
    );
};

const CartDrawer = ({ open, onClose, items, onUpdateQuantity, onRemove }: {
    open: boolean;
    onClose: () => void;
    items: CartItem[];
    onUpdateQuantity: (id: string, qty: number) => void;
    onRemove: (id: string) => void;
}) => {
    const subtotal = items.reduce((sum, item) => sum + (item.product.price * item.quantity), 0);

    return (
        <Drawer
            title={<Space><ShoppingCartOutlined /> Your Cart ({items.length})</Space>}
            open={open}
            onClose={onClose}
            width={420}
            footer={
                <div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                        <Text strong style={{ fontSize: 18 }}>Total:</Text>
                        <Text strong style={{ fontSize: 18, color: '#ff4d4f' }}>{formatNaira(subtotal)}</Text>
                    </div>
                    <Button type="primary" size="large" block icon={<CreditCardOutlined />}>
                        Checkout ({formatNaira(subtotal)})
                    </Button>
                    <Button size="large" block style={{ marginTop: 8 }}>
                        Continue Shopping
                    </Button>
                </div>
            }
        >
            {items.length === 0 ? (
                <Empty description="Your cart is empty" image={Empty.PRESENTED_IMAGE_SIMPLE}>
                    <Button type="primary" onClick={onClose}>Start Shopping</Button>
                </Empty>
            ) : (
                <List
                    itemLayout="horizontal"
                    dataSource={items}
                    renderItem={(item) => (
                        <List.Item
                            actions={[
                                <Button
                                    key="delete"
                                    type="text"
                                    danger
                                    icon={<DeleteOutlined />}
                                    onClick={() => onRemove(item.product.id)}
                                />
                            ]}
                        >
                            <List.Item.Meta
                                avatar={<Avatar size={64} shape="square" style={{ background: '#f5f5f5' }}>📦</Avatar>}
                                title={<Text ellipsis style={{ maxWidth: 180 }}>{item.product.name}</Text>}
                                description={
                                    <Space direction="vertical" size={4}>
                                        <Text strong style={{ color: '#ff4d4f' }}>{formatNaira(item.product.price)}</Text>
                                        <Space>
                                            <Button
                                                size="small"
                                                icon={<MinusOutlined />}
                                                onClick={() => onUpdateQuantity(item.product.id, Math.max(1, item.quantity - 1))}
                                            />
                                            <Text>{item.quantity}</Text>
                                            <Button
                                                size="small"
                                                icon={<PlusOutlined />}
                                                onClick={() => onUpdateQuantity(item.product.id, item.quantity + 1)}
                                            />
                                        </Space>
                                    </Space>
                                }
                            />
                        </List.Item>
                    )}
                />
            )}
        </Drawer>
    );
};

const AISearchBar = () => {
    const [isListening, setIsListening] = useState(false);

    return (
        <div style={{ position: 'relative', maxWidth: 600, flex: 1 }}>
            <Search
                placeholder="Search products, brands, categories..."
                size="large"
                enterButton={<SearchOutlined />}
                suffix={
                    <Space>
                        <Tooltip title="Voice Search">
                            <Button
                                type="text"
                                icon={<AudioOutlined style={{ color: isListening ? '#ff4d4f' : undefined }} />}
                                onClick={() => {
                                    setIsListening(!isListening);
                                    notification.info({ message: 'Voice Search', description: 'Listening... Say what you want to find!' });
                                }}
                            />
                        </Tooltip>
                        <Tooltip title="Scan Barcode">
                            <Button type="text" icon={<ScanOutlined />} />
                        </Tooltip>
                        <Tooltip title="Image Search">
                            <Button type="text" icon={<CameraOutlined />} />
                        </Tooltip>
                    </Space>
                }
            />
        </div>
    );
};

// =============================================================================
// MAIN COMPONENT
// =============================================================================

export default function ShopHomepage() {
    const [cartItems, setCartItems] = useState<CartItem[]>([]);
    const [cartOpen, setCartOpen] = useState(false);
    const [wishlistCount, setWishlistCount] = useState(0);

    const addToCart = (product: Product) => {
        setCartItems(prev => {
            const existing = prev.find(i => i.product.id === product.id);
            if (existing) {
                return prev.map(i => i.product.id === product.id ? { ...i, quantity: i.quantity + 1 } : i);
            }
            return [...prev, { product, quantity: 1 }];
        });
    };

    const updateCartQuantity = (id: string, qty: number) => {
        setCartItems(prev => prev.map(i => i.product.id === id ? { ...i, quantity: qty } : i));
    };

    const removeFromCart = (id: string) => {
        setCartItems(prev => prev.filter(i => i.product.id !== id));
    };

    const addToWishlist = (product: Product) => {
        setWishlistCount(prev => prev + 1);
    };

    return (
        <div style={{ background: '#f5f5f5', minHeight: '100vh' }}>
            {/* Promo Banner */}
            <div style={{ background: 'linear-gradient(90deg, #722ed1 0%, #1890ff 100%)', padding: '8px 24px', textAlign: 'center' }}>
                <Text style={{ color: 'white' }}>
                    🎉 <Text strong style={{ color: 'white' }}>NEW YEAR MEGA SALE!</Text> Free delivery on orders above ₦50,000.
                    Use code <Tag color="gold">OMNIROUTE2026</Tag> for extra 10% off!
                </Text>
            </div>

            {/* Header */}
            <Affix>
                <header style={{
                    background: 'white',
                    padding: '12px 24px',
                    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                }}>
                    <Row align="middle" gutter={24}>
                        <Col>
                            <Space>
                                <Avatar style={{ background: '#1a365d' }} size={40}>O</Avatar>
                                <Title level={4} style={{ margin: 0, color: '#1a365d' }}>OmniRoute</Title>
                            </Space>
                        </Col>
                        <Col flex="auto">
                            <AISearchBar />
                        </Col>
                        <Col>
                            <Space size="middle">
                                <Tooltip title="Orders">
                                    <Button type="text" icon={<HistoryOutlined style={{ fontSize: 20 }} />} />
                                </Tooltip>
                                <Tooltip title="Wishlist">
                                    <Badge count={wishlistCount}>
                                        <Button type="text" icon={<HeartOutlined style={{ fontSize: 20 }} />} />
                                    </Badge>
                                </Tooltip>
                                <Badge count={cartItems.length}>
                                    <Button
                                        type="primary"
                                        icon={<ShoppingCartOutlined />}
                                        onClick={() => setCartOpen(true)}
                                    >
                                        Cart
                                    </Button>
                                </Badge>
                                <Dropdown
                                    menu={{
                                        items: [
                                            { key: '1', icon: <UserOutlined />, label: 'My Account' },
                                            { key: '2', icon: <HistoryOutlined />, label: 'My Orders' },
                                            { key: '3', icon: <HeartOutlined />, label: 'Wishlist' },
                                            { key: '4', icon: <WalletOutlined />, label: 'Wallet' },
                                            { type: 'divider' },
                                            { key: '5', label: 'Logout' },
                                        ],
                                    }}
                                >
                                    <Avatar style={{ cursor: 'pointer', background: '#1a365d' }}>JO</Avatar>
                                </Dropdown>
                            </Space>
                        </Col>
                    </Row>

                    {/* Category Nav */}
                    <div style={{ marginTop: 12, overflowX: 'auto' }}>
                        <Space size={16}>
                            <Button type="link" icon={<AppstoreOutlined />}>All Categories</Button>
                            {categories.slice(0, 6).map(cat => (
                                <Button key={cat.id} type="text">{cat.icon} {cat.name}</Button>
                            ))}
                        </Space>
                    </div>
                </header>
            </Affix>

            {/* Hero Banner */}
            <Carousel autoplay effect="fade">
                {[
                    { title: '🔥 Flash Sale Friday!', subtitle: 'Up to 50% off on groceries', bg: 'linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%)', cta: 'Shop Now' },
                    { title: '🚚 Free Delivery', subtitle: 'On orders above ₦50,000', bg: 'linear-gradient(135deg, #52c41a 0%, #73d13d 100%)', cta: 'Learn More' },
                    { title: '💰 Bulk Discounts', subtitle: 'Save more when you buy more', bg: 'linear-gradient(135deg, #1890ff 0%, #40a9ff 100%)', cta: 'View Deals' },
                ].map((slide, i) => (
                    <div key={i}>
                        <div style={{
                            background: slide.bg,
                            padding: '48px 24px',
                            textAlign: 'center'
                        }}>
                            <Title style={{ color: 'white', marginBottom: 8, fontSize: 36 }}>{slide.title}</Title>
                            <Paragraph style={{ color: 'rgba(255,255,255,0.9)', fontSize: 18, marginBottom: 24 }}>
                                {slide.subtitle}
                            </Paragraph>
                            <Button size="large" ghost>{slide.cta} <RightOutlined /></Button>
                        </div>
                    </div>
                ))}
            </Carousel>

            {/* Value Props */}
            <div style={{ background: 'white', padding: '16px 24px', borderBottom: '1px solid #f0f0f0' }}>
                <Row gutter={24} justify="center">
                    {[
                        { icon: <TruckOutlined style={{ fontSize: 24, color: '#1890ff' }} />, text: 'Free Delivery on ₦50K+' },
                        { icon: <ThunderboltOutlined style={{ fontSize: 24, color: '#fa8c16' }} />, text: 'Same-Day Delivery' },
                        { icon: <SafetyCertificateOutlined style={{ fontSize: 24, color: '#52c41a' }} />, text: '100% Authentic Products' },
                        { icon: <CustomerServiceOutlined style={{ fontSize: 24, color: '#722ed1' }} />, text: '24/7 Customer Support' },
                    ].map((item, i) => (
                        <Col key={i}>
                            <Space>{item.icon} <Text strong>{item.text}</Text></Space>
                        </Col>
                    ))}
                </Row>
            </div>

            <div style={{ padding: '24px', maxWidth: 1400, margin: '0 auto' }}>
                {/* Categories */}
                <div style={{ marginBottom: 32 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                        <Title level={4} style={{ margin: 0 }}>Shop by Category</Title>
                        <Button type="link">View All <RightOutlined /></Button>
                    </div>
                    <Row gutter={[12, 12]}>
                        {categories.map((cat) => (
                            <Col xs={12} sm={8} md={6} lg={3} key={cat.id}>
                                <Card
                                    hoverable
                                    bodyStyle={{ padding: 16, textAlign: 'center' }}
                                    style={{ background: 'white' }}
                                >
                                    <div style={{ fontSize: 32, marginBottom: 8 }}>{cat.icon}</div>
                                    <Text strong>{cat.name}</Text>
                                    <br />
                                    <Text type="secondary" style={{ fontSize: 11 }}>{cat.productCount} products</Text>
                                </Card>
                            </Col>
                        ))}
                    </Row>
                </div>

                {/* Flash Sale Section */}
                <div style={{ marginBottom: 32 }}>
                    <Card
                        title={
                            <Space>
                                <ThunderboltOutlined style={{ color: '#ff4d4f', fontSize: 24 }} />
                                <Title level={4} style={{ margin: 0, color: '#ff4d4f' }}>Flash Sale</Title>
                                <Divider type="vertical" />
                                <Text type="secondary">Ends in:</Text>
                                <FlashSaleTimer endTime="2026-01-18T23:59:59" />
                            </Space>
                        }
                        extra={<Button type="link">View All Deals <RightOutlined /></Button>}
                        bodyStyle={{ padding: 16 }}
                    >
                        <Row gutter={[16, 16]}>
                            {featuredProducts.filter(p => p.isFlashSale || p.discount).slice(0, 4).map((product) => (
                                <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
                                    <ProductCard product={product} onAddToCart={addToCart} onAddToWishlist={addToWishlist} />
                                </Col>
                            ))}
                        </Row>
                    </Card>
                </div>

                {/* Best Sellers */}
                <div style={{ marginBottom: 32 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                        <Space>
                            <TrophyOutlined style={{ color: '#d69e2e', fontSize: 24 }} />
                            <Title level={4} style={{ margin: 0 }}>Best Sellers</Title>
                        </Space>
                        <Button type="link">View All <RightOutlined /></Button>
                    </div>
                    <Row gutter={[16, 16]}>
                        {featuredProducts.filter(p => p.isBestSeller).slice(0, 4).map((product) => (
                            <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
                                <ProductCard product={product} onAddToCart={addToCart} onAddToWishlist={addToWishlist} />
                            </Col>
                        ))}
                    </Row>
                </div>

                {/* All Products */}
                <div style={{ marginBottom: 32 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                        <Title level={4} style={{ margin: 0 }}>All Products</Title>
                        <Segmented options={['Featured', 'New Arrivals', 'Price: Low-High', 'Price: High-Low']} />
                    </div>
                    <Row gutter={[16, 16]}>
                        {featuredProducts.map((product) => (
                            <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
                                <ProductCard product={product} onAddToCart={addToCart} onAddToWishlist={addToWishlist} />
                            </Col>
                        ))}
                    </Row>
                </div>
            </div>

            {/* Footer */}
            <footer style={{ background: '#1a365d', color: 'white', padding: '48px 24px 24px' }}>
                <Row gutter={48} style={{ maxWidth: 1400, margin: '0 auto' }}>
                    <Col xs={24} md={8} style={{ marginBottom: 24 }}>
                        <Space>
                            <Avatar style={{ background: 'white', color: '#1a365d' }} size={48}>O</Avatar>
                            <Title level={3} style={{ color: 'white', margin: 0 }}>OmniRoute</Title>
                        </Space>
                        <Paragraph style={{ color: 'rgba(255,255,255,0.7)', marginTop: 16 }}>
                            Nigeria's leading B2B/B2C FMCG commerce platform. Quality products, best prices, fastest delivery.
                        </Paragraph>
                        <Space>
                            <Button type="primary" icon={<WhatsAppOutlined />} style={{ background: '#25D366', borderColor: '#25D366' }}>
                                Chat on WhatsApp
                            </Button>
                        </Space>
                    </Col>
                    <Col xs={12} md={4}>
                        <Title level={5} style={{ color: 'white' }}>Shop</Title>
                        <Space direction="vertical">
                            {['All Products', 'Categories', 'Deals', 'New Arrivals', 'Bulk Orders'].map(item => (
                                <a key={item} style={{ color: 'rgba(255,255,255,0.7)' }}>{item}</a>
                            ))}
                        </Space>
                    </Col>
                    <Col xs={12} md={4}>
                        <Title level={5} style={{ color: 'white' }}>Account</Title>
                        <Space direction="vertical">
                            {['My Orders', 'Wishlist', 'Wallet', 'Addresses', 'Settings'].map(item => (
                                <a key={item} style={{ color: 'rgba(255,255,255,0.7)' }}>{item}</a>
                            ))}
                        </Space>
                    </Col>
                    <Col xs={24} md={8}>
                        <Title level={5} style={{ color: 'white' }}>Contact Us</Title>
                        <Paragraph style={{ color: 'rgba(255,255,255,0.7)' }}>
                            📞 0800-OMNIROUTE (24/7)<br />
                            📧 support@omniroute.io<br />
                            💬 WhatsApp: +234 800 123 4567<br />
                            📍 Lagos | Abuja | Port Harcourt | Kano
                        </Paragraph>
                        <Space>
                            {['💳', '🏦', '📱'].map((icon, i) => (
                                <Tag key={i} style={{ background: 'rgba(255,255,255,0.1)', border: 'none', padding: '4px 12px' }}>{icon}</Tag>
                            ))}
                        </Space>
                    </Col>
                </Row>
                <Divider style={{ borderColor: 'rgba(255,255,255,0.1)' }} />
                <div style={{ textAlign: 'center' }}>
                    <Text style={{ color: 'rgba(255,255,255,0.5)' }}>
                        © 2026 OmniRoute Commerce Platform. All rights reserved.
                    </Text>
                </div>
            </footer>

            {/* Cart Drawer */}
            <CartDrawer
                open={cartOpen}
                onClose={() => setCartOpen(false)}
                items={cartItems}
                onUpdateQuantity={updateCartQuantity}
                onRemove={removeFromCart}
            />

            {/* Floating Actions */}
            <FloatButton.Group shape="circle" style={{ right: 24 }}>
                <FloatButton icon={<CustomerServiceOutlined />} tooltip="Chat Support" />
                <FloatButton icon={<WhatsAppOutlined style={{ color: '#25D366' }} />} tooltip="WhatsApp" />
                <FloatButton.BackTop />
            </FloatButton.Group>
        </div>
    );
}
