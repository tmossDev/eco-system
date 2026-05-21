import {
  CreateProductRequest,
  ProductDetails,
  ProductSummary,
  UpdateProductRequest,
} from './product.model';

export const MOCK_PRODUCTS: ProductSummary[] = [
  {
    id: '1',
    sku: 'GEN-MUG-001',
    name: 'Everyday Ceramic Mug',
    description: 'A durable 350ml mug for daily coffee, tea, or desk rituals.',
    category: 'Home',
    price_cents: 1299,
    currency: 'USD',
    inventory_count: 48,
    status: 'Active',
  },
  {
    id: '2',
    sku: 'APP-TEE-002',
    name: 'Organic Cotton Tee',
    description: 'Soft unisex cotton tee available in core store colors.',
    category: 'Apparel',
    price_cents: 2499,
    currency: 'USD',
    inventory_count: 82,
    status: 'Active',
  },
  {
    id: '3',
    sku: 'DIG-GUIDE-003',
    name: 'Digital Buying Guide',
    description: 'Downloadable product guide for new store customers.',
    category: 'Digital',
    price_cents: 499,
    currency: 'USD',
    inventory_count: 999,
    status: 'Draft',
  },
  {
    id: '4',
    sku: 'KIT-STARTER-004',
    name: 'Starter Gift Kit',
    description: 'A bundled kit that can represent arbitrary grouped products.',
    category: 'Bundles',
    price_cents: 5499,
    currency: 'USD',
    inventory_count: 16,
    status: 'Archived',
  },
];

let mockProducts = [...MOCK_PRODUCTS];

export function getMockProducts(): ProductSummary[] {
  return [...mockProducts];
}

export function getMockProductById(id: string): ProductDetails | undefined {
  return mockProducts.find((product) => product.id === id);
}

export function createMockProduct(request: CreateProductRequest): ProductDetails {
  const product: ProductDetails = {
    id: createMockProductId(),
    ...request,
  };

  mockProducts = [product, ...mockProducts];

  return product;
}

export function updateMockProduct(
  id: string,
  request: UpdateProductRequest,
): ProductDetails {
  const updatedProduct: ProductDetails = {
    id,
    ...request,
  };

  mockProducts = mockProducts.map((product) =>
    product.id === id ? updatedProduct : product,
  );

  return updatedProduct;
}

export function deleteMockProduct(id: string): void {
  mockProducts = mockProducts.filter((product) => product.id !== id);
}

export function createMockProductId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID();
  }

  return String(Date.now());
}
