import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HostreviewUpdateComponent } from './hostreview-update.component';

describe('HostreviewUpdateComponent', () => {
  let component: HostreviewUpdateComponent;
  let fixture: ComponentFixture<HostreviewUpdateComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [HostreviewUpdateComponent]
    });
    fixture = TestBed.createComponent(HostreviewUpdateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
