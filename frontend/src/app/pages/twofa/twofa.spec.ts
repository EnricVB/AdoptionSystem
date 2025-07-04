import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Twofa } from './twofa';

describe('Twofa', () => {
  let component: Twofa;
  let fixture: ComponentFixture<Twofa>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Twofa]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Twofa);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
