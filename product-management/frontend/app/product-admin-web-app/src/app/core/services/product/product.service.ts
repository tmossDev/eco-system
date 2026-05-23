import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {
  CreateDiscountRequest,
  CreateProductRequest,
  Discount,
  ProductDetails,
  ProductSummary,
  UpdateDiscountRequest,
  UpdateProductRequest,
} from './product.model';

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/products';

  public getProducts(): Observable<ProductSummary[]> {
    return this.http.get<ProductSummary[]>(this.baseUrl);
  }

  public getProductById(id: string): Observable<ProductDetails> {
    return this.http.get<ProductDetails>(`${this.baseUrl}/${id}`);
  }

  public createProduct(request: CreateProductRequest): Observable<ProductDetails> {
    return this.http.post<ProductDetails>(this.baseUrl, request);
  }

  public updateProduct(
    id: string,
    request: UpdateProductRequest,
  ): Observable<ProductDetails> {
    return this.http.put<ProductDetails>(`${this.baseUrl}/${id}`, request);
  }

  public uploadProductPhoto(id: string, file: File): Observable<ProductDetails> {
    const formData = new FormData();
    formData.append('file', file);

    return this.http.post<ProductDetails>(
      `${this.baseUrl}/${id}/photos`,
      formData,
    );
  }

  public deleteProduct(id: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }

  public getDiscounts(): Observable<Discount[]> {
    return this.http.get<Discount[]>('/api/discounts');
  }

  public getDiscountById(id: string): Observable<Discount> {
    return this.http.get<Discount>(`/api/discounts/${id}`);
  }

  public createDiscount(request: CreateDiscountRequest): Observable<Discount> {
    return this.http.post<Discount>('/api/discounts', request);
  }

  public updateDiscount(
    id: string,
    request: UpdateDiscountRequest,
  ): Observable<Discount> {
    return this.http.put<Discount>(`/api/discounts/${id}`, request);
  }

  public deleteDiscount(id: string): Observable<void> {
    return this.http.delete<void>(`/api/discounts/${id}`);
  }
}
