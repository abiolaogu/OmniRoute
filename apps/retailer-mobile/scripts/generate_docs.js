const {
  Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell,
  Header, Footer, AlignmentType, HeadingLevel, BorderStyle, WidthType,
  ShadingType, PageNumber, PageBreak, LevelFormat
} = require('docx');
const fs = require('fs');

// Constants
const border = { style: BorderStyle.SINGLE, size: 1, color: "CCCCCC" };
const borders = { top: border, bottom: border, left: border, right: border };
const cellMargins = { top: 80, bottom: 80, left: 120, right: 120 };

// Helper functions
function createHeading(text, level) {
  return new Paragraph({
    heading: level,
    children: [new TextRun({ text, bold: true })]
  });
}

function createParagraph(text, options = {}) {
  return new Paragraph({
    children: [new TextRun({ text, ...options })]
  });
}

function createBulletList(items, reference) {
  return items.map(item => new Paragraph({
    numbering: { reference, level: 0 },
    children: [new TextRun(item)]
  }));
}

// Create document
const doc = new Document({
  styles: {
    default: {
      document: { run: { font: "Arial", size: 24 } }
    },
    paragraphStyles: [
      {
        id: "Heading1", name: "Heading 1", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 36, bold: true, font: "Arial", color: "0D47A1" },
        paragraph: { spacing: { before: 360, after: 240 }, outlineLevel: 0 }
      },
      {
        id: "Heading2", name: "Heading 2", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 28, bold: true, font: "Arial", color: "1565C0" },
        paragraph: { spacing: { before: 240, after: 180 }, outlineLevel: 1 }
      },
      {
        id: "Heading3", name: "Heading 3", basedOn: "Normal", next: "Normal", quickFormat: true,
        run: { size: 24, bold: true, font: "Arial", color: "1976D2" },
        paragraph: { spacing: { before: 180, after: 120 }, outlineLevel: 2 }
      }
    ]
  },
  numbering: {
    config: [
      {
        reference: "bullets",
        levels: [{
          level: 0, format: LevelFormat.BULLET, text: "•", alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 720, hanging: 360 } } }
        }]
      },
      {
        reference: "numbers",
        levels: [{
          level: 0, format: LevelFormat.DECIMAL, text: "%1.", alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 720, hanging: 360 } } }
        }]
      }
    ]
  },
  sections: [{
    properties: {
      page: {
        size: { width: 12240, height: 15840 },
        margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 }
      }
    },
    headers: {
      default: new Header({
        children: [new Paragraph({
          alignment: AlignmentType.RIGHT,
          children: [
            new TextRun({ text: "OmniRoute Ecosystem - Technical Documentation", color: "666666", size: 20 })
          ]
        })]
      })
    },
    footers: {
      default: new Footer({
        children: [new Paragraph({
          alignment: AlignmentType.CENTER,
          children: [
            new TextRun({ text: "Page ", size: 20 }),
            new TextRun({ children: [PageNumber.CURRENT], size: 20 }),
            new TextRun({ text: " of ", size: 20 }),
            new TextRun({ children: [PageNumber.TOTAL_PAGES], size: 20 })
          ]
        })]
      })
    },
    children: [
      // Title Page
      new Paragraph({ spacing: { before: 2000 } }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "OmniRoute Ecosystem", size: 72, bold: true, color: "0D47A1" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 200, after: 400 },
        children: [new TextRun({ text: "Technical Specification Document", size: 36, color: "666666" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 400 },
        children: [new TextRun({ text: "B2B FMCG Multi-Participant Commerce Platform", size: 28, italics: true })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 800 },
        children: [new TextRun({ text: "Version 1.0", size: 24 })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 200 },
        children: [new TextRun({ text: "January 2026", size: 24, color: "666666" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 1000 },
        children: [new TextRun({ text: "BillyRonks Global Limited", size: 24, bold: true })]
      }),

      // Page break
      new Paragraph({ children: [new PageBreak()] }),

      // Executive Summary
      createHeading("Executive Summary", HeadingLevel.HEADING_1),
      createParagraph("OmniRoute Ecosystem is a comprehensive mobile-first B2B FMCG platform built with Flutter that provides tailored experiences for all participants in the commerce ecosystem. The platform connects banks, financial institutions, logistics companies, warehouses, manufacturers, e-commerce businesses, retailers, wholesalers, entrepreneurs, investors, agents, and drivers."),
      createParagraph(""),
      createParagraph("This document provides technical specifications for the Flutter mobile application, covering architecture, design patterns, implementation details, and development guidelines."),
      createParagraph(""),

      // Key Metrics
      createHeading("Project Metrics", HeadingLevel.HEADING_2),
      new Table({
        width: { size: 100, type: WidthType.PERCENTAGE },
        columnWidths: [4680, 4680],
        rows: [
          new TableRow({
            children: [
              new TableCell({
                borders, width: { size: 4680, type: WidthType.DXA },
                shading: { fill: "E3F2FD", type: ShadingType.CLEAR },
                margins: cellMargins,
                children: [new Paragraph({ children: [new TextRun({ text: "Metric", bold: true })] })]
              }),
              new TableCell({
                borders, width: { size: 4680, type: WidthType.DXA },
                shading: { fill: "E3F2FD", type: ShadingType.CLEAR },
                margins: cellMargins,
                children: [new Paragraph({ children: [new TextRun({ text: "Value", bold: true })] })]
              })
            ]
          }),
          ...[ 
            ["Total Dart Files", "27"],
            ["Lines of Code", "13,000+"],
            ["Participant Types", "12"],
            ["Dashboard Variants", "10"],
            ["Feature Modules", "8"],
            ["Framework", "Flutter 3.2+"]
          ].map(([metric, value]) => new TableRow({
            children: [
              new TableCell({ borders, width: { size: 4680, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(metric)] }),
              new TableCell({ borders, width: { size: 4680, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(value)] })
            ]
          }))
        ]
      }),

      new Paragraph({ children: [new PageBreak()] }),

      // Architecture
      createHeading("System Architecture", HeadingLevel.HEADING_1),

      createHeading("Technology Stack", HeadingLevel.HEADING_2),
      new Table({
        width: { size: 100, type: WidthType.PERCENTAGE },
        columnWidths: [3120, 3120, 3120],
        rows: [
          new TableRow({
            children: [
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, shading: { fill: "0D47A1", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Layer", bold: true, color: "FFFFFF" })] })] }),
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, shading: { fill: "0D47A1", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Technology", bold: true, color: "FFFFFF" })] })] }),
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, shading: { fill: "0D47A1", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Purpose", bold: true, color: "FFFFFF" })] })] })
            ]
          }),
          ...[ 
            ["Framework", "Flutter 3.2+", "Cross-platform UI"],
            ["State Management", "Riverpod", "Reactive state"],
            ["Navigation", "GoRouter", "Declarative routing"],
            ["Network", "Dio", "HTTP client"],
            ["Local Storage", "Hive", "Offline data"],
            ["Security", "FlutterSecureStorage", "Token storage"],
            ["Charts", "FL Chart", "Data visualization"],
            ["Animation", "flutter_animate", "UI animations"]
          ].map(([layer, tech, purpose]) => new TableRow({
            children: [
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(layer)] }),
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(tech)] }),
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(purpose)] })
            ]
          }))
        ]
      }),

      createParagraph(""),

      createHeading("Project Structure", HeadingLevel.HEADING_2),
      createParagraph("The application follows a feature-first architecture with clear separation of concerns:"),
      createParagraph(""),
      ...createBulletList([
        "core/ - Foundation layer (constants, theme, network, routing)",
        "features/ - Feature modules (auth, dashboard, orders, inventory, wallet, settings)",
        "models/ - Data classes using Freezed for immutability",
        "providers/ - Riverpod state management providers",
        "widgets/ - Reusable UI components"
      ], "bullets"),

      new Paragraph({ children: [new PageBreak()] }),

      // Participant Types
      createHeading("Participant Types", HeadingLevel.HEADING_1),
      createParagraph("The platform supports 12 distinct participant types, each with tailored dashboards and functionality:"),
      createParagraph(""),

      new Table({
        width: { size: 100, type: WidthType.PERCENTAGE },
        columnWidths: [2340, 4680, 2340],
        rows: [
          new TableRow({
            children: [
              new TableCell({ borders, width: { size: 2340, type: WidthType.DXA }, shading: { fill: "0D47A1", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Type", bold: true, color: "FFFFFF" })] })] }),
              new TableCell({ borders, width: { size: 4680, type: WidthType.DXA }, shading: { fill: "0D47A1", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Key Features", bold: true, color: "FFFFFF" })] })] }),
              new TableCell({ borders, width: { size: 2340, type: WidthType.DXA }, shading: { fill: "0D47A1", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Color", bold: true, color: "FFFFFF" })] })] })
            ]
          }),
          ...[ 
            ["Bank", "Loans, settlements, ATC (Authority to Collect), compliance", "#1565C0"],
            ["Logistics", "Fleet management, delivery tracking, route optimization", "#E65100"],
            ["Warehouse", "Inventory management, inbound/outbound operations", "#6A1B9A"],
            ["Manufacturer", "Product catalog, order fulfillment, production tracking", "#2E7D32"],
            ["Distributor", "Distribution management, territory allocation", "#00838F"],
            ["Wholesaler", "Bulk orders, retailer management, pricing", "#4527A0"],
            ["Retailer", "POS integration, inventory, reordering, BNPL", "#C62828"],
            ["E-commerce", "Dropshipping, marketplace integrations, multi-channel", "#AD1457"],
            ["Entrepreneur", "Business ideas, learning resources, networking", "#558B2F"],
            ["Investor", "Portfolio management, investment opportunities", "#283593"],
            ["Agent", "Field tasks, commission tracking, performance", "#00695C"],
            ["Driver", "Deliveries, earnings, real-time navigation", "#BF360C"]
          ].map(([type, features, color]) => new TableRow({
            children: [
              new TableCell({ borders, width: { size: 2340, type: WidthType.DXA }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: type, bold: true })] })] }),
              new TableCell({ borders, width: { size: 4680, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(features)] }),
              new TableCell({ borders, width: { size: 2340, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(color)] })
            ]
          }))
        ]
      }),

      new Paragraph({ children: [new PageBreak()] }),

      // Design System
      createHeading("Design System", HeadingLevel.HEADING_1),

      createHeading("Color Palette", HeadingLevel.HEADING_2),
      createParagraph("The application uses Material Design 3 with custom brand colors:"),
      createParagraph(""),
      ...createBulletList([
        "Primary: #0D47A1 (Deep Blue) - Main brand color",
        "Secondary: #00BFA5 (Teal) - Secondary actions",
        "Accent: #FF6D00 (Orange) - Highlights and CTAs",
        "Success: #2E7D32 (Green) - Positive states",
        "Warning: #F57C00 (Amber) - Caution states",
        "Error: #C62828 (Red) - Error states"
      ], "bullets"),

      createHeading("Typography", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Display/Headings: Space Grotesk - Modern geometric sans-serif",
        "Body Text: Inter - Highly readable, variable font",
        "Responsive scaling with min/max constraints (0.8x - 1.2x)"
      ], "bullets"),

      createHeading("Reusable Components", HeadingLevel.HEADING_2),
      new Table({
        width: { size: 100, type: WidthType.PERCENTAGE },
        columnWidths: [3120, 6240],
        rows: [
          new TableRow({
            children: [
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, shading: { fill: "E3F2FD", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Component", bold: true })] })] }),
              new TableCell({ borders, width: { size: 6240, type: WidthType.DXA }, shading: { fill: "E3F2FD", type: ShadingType.CLEAR }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: "Purpose", bold: true })] })] })
            ]
          }),
          ...[ 
            ["StatCard", "Dashboard statistics with icon, value, and growth indicator"],
            ["WalletCard", "Balance display with gradient background and action buttons"],
            ["OrderListTile", "Order summary with status chip and formatted currency"],
            ["DeliveryListTile", "Delivery tracking with driver info and status"],
            ["StatusChip", "Color-coded status indicators (pending, processing, completed)"],
            ["QuickActionGrid", "Grid of quick action buttons with optional badges"],
            ["SectionHeader", "Section title with optional action button"],
            ["EmptyState", "Empty state with icon, message, and action"],
            ["ShimmerLoading", "Animated loading placeholder"]
          ].map(([component, purpose]) => new TableRow({
            children: [
              new TableCell({ borders, width: { size: 3120, type: WidthType.DXA }, margins: cellMargins, children: [new Paragraph({ children: [new TextRun({ text: component, italics: true })] })] }),
              new TableCell({ borders, width: { size: 6240, type: WidthType.DXA }, margins: cellMargins, children: [createParagraph(purpose)] })
            ]
          }))
        ]
      }),

      new Paragraph({ children: [new PageBreak()] }),

      // State Management
      createHeading("State Management", HeadingLevel.HEADING_1),
      createParagraph("The application uses Riverpod for reactive state management with clear separation between UI and business logic."),
      createParagraph(""),

      createHeading("Core Providers", HeadingLevel.HEADING_2),
      ...createBulletList([
        "authProvider - Authentication state (login, register, token management)",
        "ordersProvider - Paginated orders with filtering",
        "dashboardStatsProvider - Real-time dashboard statistics",
        "inventoryProvider - Stock data with alerts",
        "walletProvider - Balance and transaction history",
        "notificationsProvider - Push notification state"
      ], "bullets"),

      createHeading("Benefits of Riverpod", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Compile-time safety with type checking",
        "Automatic disposal of unused providers",
        "Easy testing with provider overrides",
        "Efficient rebuilds with selective listening",
        "Support for async operations (FutureProvider, StreamProvider)"
      ], "bullets"),

      new Paragraph({ children: [new PageBreak()] }),

      // Navigation
      createHeading("Navigation Architecture", HeadingLevel.HEADING_1),
      createParagraph("GoRouter provides declarative navigation with authentication guards and deep linking support."),
      createParagraph(""),

      createHeading("Route Types", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Public routes: splash, welcome, login, register, OTP verification",
        "Protected routes: dashboard, orders, inventory, wallet, settings",
        "Shell routes: persistent bottom navigation wrapper"
      ], "bullets"),

      createHeading("Participant-Specific Navigation", HeadingLevel.HEADING_2),
      createParagraph("Each participant type has customized bottom navigation items:"),
      createParagraph(""),
      ...createBulletList([
        "Bank: Home, Loans, Settlements, Analytics, Profile",
        "Logistics: Home, Deliveries, Fleet, Routes, Profile",
        "Retailer: Home, Orders, Inventory, Wallet, Profile",
        "E-commerce: Home, Products, Orders, Shipping, Profile"
      ], "bullets"),

      new Paragraph({ children: [new PageBreak()] }),

      // API Integration
      createHeading("API Integration", HeadingLevel.HEADING_1),

      createHeading("API Client Features", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Automatic Bearer token injection from secure storage",
        "Token refresh on 401 with automatic request retry",
        "Retry interceptor for network failures (max 3 retries)",
        "Typed responses with ApiResponse<T> wrapper",
        "File upload with progress tracking",
        "Comprehensive error message extraction",
        "Debug logging in development mode"
      ], "bullets"),

      createHeading("Data Models", HeadingLevel.HEADING_2),
      createParagraph("All models use Freezed for immutability and JSON serialization:"),
      createParagraph(""),
      ...createBulletList([
        "User - Profile with participant type and verification status",
        "Order - Complete order data with items, totals, and status",
        "Product - Product catalog with variants and pricing",
        "Delivery - Tracking with coordinates and status updates",
        "Wallet - Balance, pending amounts, and currency",
        "Transaction - Payment history with metadata"
      ], "bullets"),

      new Paragraph({ children: [new PageBreak()] }),

      // Security
      createHeading("Security Implementation", HeadingLevel.HEADING_1),
      ...createBulletList([
        "FlutterSecureStorage for sensitive token storage",
        "Automatic token refresh before expiration",
        "Biometric authentication (fingerprint/Face ID)",
        "SSL certificate pinning (production configuration)",
        "Input validation on all forms",
        "Secure API communication with HTTPS"
      ], "bullets"),

      createHeading("Authentication Flow", HeadingLevel.HEADING_2),
      new Paragraph({
        numbering: { reference: "numbers", level: 0 },
        children: [new TextRun("User enters credentials (email/phone + password)")]
      }),
      new Paragraph({
        numbering: { reference: "numbers", level: 0 },
        children: [new TextRun("API validates and returns access + refresh tokens")]
      }),
      new Paragraph({
        numbering: { reference: "numbers", level: 0 },
        children: [new TextRun("Tokens stored securely in FlutterSecureStorage")]
      }),
      new Paragraph({
        numbering: { reference: "numbers", level: 0 },
        children: [new TextRun("API client auto-injects Bearer token on requests")]
      }),
      new Paragraph({
        numbering: { reference: "numbers", level: 0 },
        children: [new TextRun("On 401, refresh token used to obtain new access token")]
      }),
      new Paragraph({
        numbering: { reference: "numbers", level: 0 },
        children: [new TextRun("Failed refresh triggers logout and redirect to login")]
      }),

      new Paragraph({ children: [new PageBreak()] }),

      // Development
      createHeading("Development Guidelines", HeadingLevel.HEADING_1),

      createHeading("Code Organization", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Feature-first folder structure for scalability",
        "Single responsibility principle for files",
        "Consistent naming conventions (snake_case for files, PascalCase for classes)",
        "Separation of UI, state, and business logic"
      ], "bullets"),

      createHeading("Testing Strategy", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Unit tests for providers and utilities",
        "Widget tests for UI components",
        "Integration tests for critical user flows",
        "Provider mocking for isolated testing"
      ], "bullets"),

      createHeading("Performance Optimization", HeadingLevel.HEADING_2),
      ...createBulletList([
        "Lazy loading with pagination",
        "Image caching with cached_network_image",
        "Shimmer loading states for better perceived performance",
        "Selective state rebuilds with Riverpod select()",
        "const constructors for static widgets"
      ], "bullets"),

      new Paragraph({ children: [new PageBreak()] }),

      // Conclusion
      createHeading("Conclusion", HeadingLevel.HEADING_1),
      createParagraph("The OmniRoute Ecosystem Flutter application provides a robust, scalable foundation for the B2B FMCG multi-participant platform. The architecture emphasizes:"),
      createParagraph(""),
      ...createBulletList([
        "Maintainability through clear separation of concerns",
        "Scalability with feature-first modular architecture",
        "Performance with optimized state management",
        "Security with proper token handling and storage",
        "User experience with consistent design system"
      ], "bullets"),
      createParagraph(""),
      createParagraph("The platform is ready for integration with backend services and deployment to app stores following the standard Flutter build and release processes."),

      createParagraph(""),
      createParagraph(""),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 400 },
        children: [new TextRun({ text: "— End of Document —", color: "666666", italics: true })]
      })
    ]
  }]
});

// Generate document
Packer.toBuffer(doc).then(buffer => {
  fs.writeFileSync('/mnt/user-data/outputs/OmniRoute_Ecosystem_Technical_Specification.docx', buffer);
  console.log('Document created successfully!');
}).catch(err => {
  console.error('Error creating document:', err);
});
