import { Discount } from './product.model';

export interface DiscountedPrice {
  originalCents: number;
  finalCents: number;
  discountCents: number;
  hasDiscount: boolean;
}

export function activeDiscounts(discounts: Discount[] | null | undefined): Discount[] {
  const now = Date.now();

  return (discounts ?? []).filter((discount) => {
    const startsAt = Date.parse(discount.starts_at);
    const endsAt = Date.parse(discount.ends_at);

    return (
      discount.status === 'Active' &&
      (Number.isNaN(startsAt) || startsAt <= now) &&
      (Number.isNaN(endsAt) || endsAt >= now)
    );
  });
}

export function calculateDiscountedPrice(
  priceCents: number,
  discounts: Discount[] | null | undefined,
): DiscountedPrice {
  const discountCents = activeDiscounts(discounts).reduce(
    (total, discount) => total + calculateDiscountCents(priceCents, discount),
    0,
  );
  const clampedDiscount = Math.min(priceCents, Math.max(0, discountCents));
  const finalCents = priceCents - clampedDiscount;

  return {
    originalCents: priceCents,
    finalCents,
    discountCents: clampedDiscount,
    hasDiscount: clampedDiscount > 0,
  };
}

export function formatDiscountValue(
  discount: Pick<Discount, 'discount_type' | 'percentage_basis_points' | 'amount_cents' | 'currency'>,
  formatMoney: (priceCents: number, currency: string) => string,
): string {
  if (discount.discount_type === 'Percentage') {
    return `${(discount.percentage_basis_points ?? 0) / 100}%`;
  }

  return formatMoney(discount.amount_cents ?? 0, discount.currency || 'USD');
}

function calculateDiscountCents(priceCents: number, discount: Discount): number {
  if (discount.discount_type === 'Percentage') {
    return Math.round((priceCents * (discount.percentage_basis_points ?? 0)) / 10000);
  }

  return discount.amount_cents ?? 0;
}
