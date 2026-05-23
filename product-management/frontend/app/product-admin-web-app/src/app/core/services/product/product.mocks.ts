import {
  CreateDiscountRequest,
  CreateProductRequest,
  Discount,
  ProductDetails,
  ProductSummary,
  UpdateDiscountRequest,
  UpdateProductRequest,
} from './product.model';

export const MOCK_DISCOUNTS: Discount[] = [
  {
    id: '1',
    name: 'Store-wide launch offer',
    description: 'Applies to every active product in the catalog.',
    discount_type: 'Percentage',
    scope: 'Global',
    percentage_basis_points: 1000,
    amount_cents: null,
    currency: '',
    min_product_count: 1,
    starts_at: '',
    ends_at: '',
    status: 'Active',
    product_ids: [],
  },
  {
    id: '2',
    name: 'Starter bundle savings',
    description: 'Qualifies when the selected kit products are purchased together.',
    discount_type: 'Amount',
    scope: 'ProductSet',
    percentage_basis_points: null,
    amount_cents: 500,
    currency: 'USD',
    min_product_count: 2,
    starts_at: '',
    ends_at: '',
    status: 'Draft',
    product_ids: ['1', '4'],
  },
];

export const MOCK_PRODUCTS: ProductSummary[] = [
  {
    id: '1',
    sku: 'GEN-MUG-001',
    name: 'Everyday Ceramic Mug',
    short_description: 'Durable 350ml ceramic mug for daily coffee and tea.',
    description: 'A durable 350ml mug for daily coffee, tea, or desk rituals.',
    category: 'Home',
    price_cents: 1299,
    currency: 'USD',
    inventory_count: 48,
    status: 'Active',
    discounts: [],
    labels: ['coffee', 'home', 'giftable'],
    photos: [
      {
        url: 'https://images.unsplash.com/photo-1514228742587-6b1558fcca3d?auto=format&fit=crop&w=1200&q=80',
        thumbnail_url:
          'https://images.unsplash.com/photo-1514228742587-6b1558fcca3d?auto=format&fit=crop&w=160&q=70',
        alt_text: 'White ceramic mug on a table',
        is_primary: true,
      },
    ],
  },
  {
    id: '2',
    sku: 'APP-TEE-002',
    name: 'Organic Cotton Tee',
    short_description: 'Soft unisex tee in core store colors.',
    description: 'Soft unisex cotton tee available in core store colors.',
    category: 'Apparel',
    price_cents: 2499,
    currency: 'USD',
    inventory_count: 82,
    status: 'Active',
    discounts: [],
    labels: ['apparel', 'organic'],
    photos: [
      {
        url: 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=1200&q=80',
        thumbnail_url:
          'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=160&q=70',
        alt_text: 'Folded cotton tee',
        is_primary: true,
      },
    ],
  },
  {
    id: '3',
    sku: 'DIG-GUIDE-003',
    name: 'Digital Buying Guide',
    short_description: 'Downloadable buying guide for new store customers.',
    description: 'Downloadable product guide for new store customers.',
    category: 'Digital',
    price_cents: 499,
    currency: 'USD',
    inventory_count: 999,
    status: 'Draft',
    discounts: [],
    labels: ['digital', 'guide'],
    photos: [],
  },
  {
    id: '4',
    sku: 'KIT-STARTER-004',
    name: 'Starter Gift Kit',
    short_description: 'Bundle-ready starter kit for curated product sets.',
    description: 'A bundled kit that can represent arbitrary grouped products.',
    category: 'Bundles',
    price_cents: 5499,
    currency: 'USD',
    inventory_count: 16,
    status: 'Archived',
    discounts: [],
    labels: ['bundle', 'giftable'],
    photos: [
      {
        url: 'https://images.unsplash.com/photo-1549465220-1a8b9238cd48?auto=format&fit=crop&w=1200&q=80',
        thumbnail_url:
          'https://images.unsplash.com/photo-1549465220-1a8b9238cd48?auto=format&fit=crop&w=160&q=70',
        alt_text: 'Wrapped gift kit box',
        is_primary: true,
      },
    ],
  },
];

let mockProducts = [...MOCK_PRODUCTS];
let mockDiscounts = [...MOCK_DISCOUNTS];

export function getMockProducts(): ProductSummary[] {
  return mockProducts.map(withApplicableDiscounts);
}

export function getMockProductById(id: string): ProductDetails | undefined {
  const product = mockProducts.find((product) => product.id === id);

  return product ? withApplicableDiscounts(product) : undefined;
}

export function createMockProduct(request: CreateProductRequest): ProductDetails {
  const product: ProductDetails = {
    id: createMockProductId(),
    ...request,
    labels: request.labels ?? [],
    discounts: [],
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
    labels: request.labels ?? [],
    discounts: [],
  };

  mockProducts = mockProducts.map((product) =>
    product.id === id ? updatedProduct : product,
  );

  return updatedProduct;
}

export function uploadMockProductPhoto(id: string, fileName: string): ProductDetails {
  const product = getMockProductById(id);
  if (!product) {
    throw new Error(`Product with id "${id}" was not found`);
  }

  const updatedProduct: ProductDetails = {
    ...product,
    photos: [
      ...product.photos,
      {
        url: '/assets/favicon.ico',
        thumbnail_url: '/assets/favicon.ico',
        alt_text: fileName || product.name,
        is_primary: product.photos.length === 0,
      },
    ],
  };

  mockProducts = mockProducts.map((candidate) =>
    candidate.id === id ? updatedProduct : candidate,
  );

  return updatedProduct;
}

export function deleteMockProduct(id: string): void {
  mockProducts = mockProducts.filter((product) => product.id !== id);
}

export function getMockDiscounts(): Discount[] {
  return [...mockDiscounts];
}

export function getMockDiscountById(id: string): Discount | undefined {
  return mockDiscounts.find((discount) => discount.id === id);
}

export function createMockDiscount(request: CreateDiscountRequest): Discount {
  const discount: Discount = {
    id: createMockProductId(),
    ...request,
  };

  mockDiscounts = [discount, ...mockDiscounts];

  return discount;
}

export function updateMockDiscount(
  id: string,
  request: UpdateDiscountRequest,
): Discount {
  const updatedDiscount: Discount = {
    id,
    ...request,
  };

  mockDiscounts = mockDiscounts.map((discount) =>
    discount.id === id ? updatedDiscount : discount,
  );

  return updatedDiscount;
}

export function deleteMockDiscount(id: string): void {
  mockDiscounts = mockDiscounts.filter((discount) => discount.id !== id);
}

export function createMockProductId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID();
  }

  return String(Date.now());
}

function withApplicableDiscounts<Product extends ProductSummary>(
  product: Product,
): Product {
  const discounts = mockDiscounts.filter(
    (discount) =>
      discount.scope === 'Global' ||
      discount.product_ids.some((productID) => String(productID) === String(product.id)),
  );

  return {
    ...product,
    discounts,
  };
}
