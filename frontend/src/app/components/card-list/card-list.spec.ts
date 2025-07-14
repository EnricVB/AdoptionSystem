import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CardList } from './card-list';

describe('CardList', () => {
  let component: CardList;
  let fixture: ComponentFixture<CardList>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CardList]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CardList);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
