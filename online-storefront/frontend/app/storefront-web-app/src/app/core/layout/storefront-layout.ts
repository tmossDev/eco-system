import { Component, inject, signal } from '@angular/core';
import { Router, RouterLink, RouterOutlet } from '@angular/router';
import { finalize } from 'rxjs';

import { AuthService } from '../services/auth.service';
import { CartService } from '../services/cart.service';

@Component({
  selector: 'app-storefront-layout',
  imports: [RouterLink, RouterOutlet],
  template: `
    <header>
      <a class="brand" routerLink="/" aria-label="Northstar Store home">
        <span>N</span>
        <strong>Northstar</strong>
      </a>

      <nav aria-label="Main navigation">
        <a routerLink="/">Shop</a>
        <a href="#about">Our story</a>
      </nav>

      <div class="account">
        <a class="cart" routerLink="/cart">
          Cart
          @if (cartService.itemCount()) {
            <span>{{ cartService.itemCount() }}</span>
          }
        </a>
        @if (authService.isAuthenticated()) {
          <span class="welcome">Hi, {{ firstName() }}</span>
          <button type="button" [disabled]="isSigningOut()" (click)="logout()">
            {{ isSigningOut() ? 'Signing out...' : 'Sign out' }}
          </button>
        } @else {
          <a class="sign-in" routerLink="/auth/login">Sign in</a>
          <a class="join" routerLink="/auth/register">Create account</a>
        }
      </div>
    </header>

    <main>
      <router-outlet />
    </main>

    <footer id="about">
      <div>
        <a class="brand" routerLink="/"><span>N</span><strong>Northstar</strong></a>
        <p>Useful things, carefully chosen for everyday life.</p>
      </div>
      <p>Independent catalog. Simple shopping.</p>
    </footer>
  `,
  styles: `
    header, footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1.5rem;
      padding: 1.1rem clamp(1rem, 5vw, 5rem);
    }
    header { border-bottom: 1px solid #dfded4; background: rgb(255 253 248 / 92%); }
    .brand { display: inline-flex; align-items: center; gap: .6rem; text-decoration: none; }
    .brand span { display: grid; width: 2rem; height: 2rem; place-items: center; border-radius: 50%; background: #31543c; color: white; font-family: Georgia, serif; }
    .brand strong { font-family: Georgia, serif; font-size: 1.35rem; }
    nav { display: flex; gap: 1.5rem; margin-right: auto; }
    nav a, .sign-in, .cart { color: #566459; font-weight: 700; text-decoration: none; }
    .account { display: flex; align-items: center; gap: .85rem; }
    .cart { display: inline-flex; align-items: center; gap: .35rem; }
    .cart span { display: grid; min-width: 1.35rem; height: 1.35rem; place-items: center; border-radius: 999px; background: #31543c; color: white; font-size: .72rem; }
    .welcome { color: #566459; font-size: .9rem; }
    button, .join { border: 0; border-radius: 999px; background: #31543c; padding: .7rem 1rem; color: white; cursor: pointer; font-weight: 700; text-decoration: none; }
    button:disabled { cursor: wait; opacity: .7; }
    footer { align-items: flex-start; margin-top: 4rem; border-top: 1px solid #dfded4; background: #ede8dc; color: #657067; }
    footer p { margin: .65rem 0 0; font-size: .9rem; }
    @media (max-width: 720px) {
      header { flex-wrap: wrap; }
      nav { order: 3; width: 100%; }
      .welcome { display: none; }
      footer { display: block; }
    }
  `,
})
export class StorefrontLayout {
  protected readonly authService = inject(AuthService);
  protected readonly cartService = inject(CartService);
  private readonly router = inject(Router);
  protected readonly isSigningOut = signal(false);

  protected firstName(): string {
    return this.authService.currentUser()?.name.split(' ')[0] ?? 'there';
  }

  protected logout(): void {
    this.isSigningOut.set(true);
    this.authService
      .logout()
      .pipe(finalize(() => this.isSigningOut.set(false)))
      .subscribe({
        next: () => void this.router.navigate(['/']),
        error: () => void this.router.navigate(['/']),
      });
  }
}
