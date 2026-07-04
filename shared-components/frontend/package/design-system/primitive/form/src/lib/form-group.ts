import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

let nextGroupId = 0;

@Component({
  selector: 'primitive-form-group',
  standalone: true,
  template: `
    @if (legend) {
      <legend class="primitive-form-group__legend" [id]="legendId">{{ legend }}</legend>
    }

    @if (description) {
      <p class="primitive-form-group__description" [id]="descriptionId">{{ description }}</p>
    }

    <div class="primitive-form-group__content">
      <ng-content />
    </div>
  `,
  styleUrl: './form.scss',
  host: {
    class: 'primitive-form-group',
    '[attr.role]': '"group"',
    '[attr.aria-labelledby]': 'legend ? legendId : null',
    '[attr.aria-describedby]': 'description ? descriptionId : null',
  },
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DsFormGroup {
  @Input()
  public legend = '';

  @Input()
  public description = '';

  protected readonly groupId = `primitive-form-group-${nextGroupId++}`;
  protected readonly legendId = `${this.groupId}-legend`;
  protected readonly descriptionId = `${this.groupId}-description`;
}
