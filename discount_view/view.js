function initDiscountView() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        loading: true,
        discount: null
      };
    },
    mounted() {
      this.loadDiscountDetails();
    },
    methods: {
      async loadDiscountDetails() {
        this.loading = true;
        try {
          const response = await fetch(urlLoadDetails, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              discount_id: DISCOUNT_ID
            })
          });
          const data = await response.json();

          if (data.status === 'success') {
            this.discount = data.data || null;
          } else {
            Notiflix.Notify.failure(data.message || 'Failed to load discount details', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        } catch (error) {
          console.error('Error loading discount details:', error);
          Notiflix.Notify.failure(error.message || 'Failed to load discount details', {
            position: 'right-top',
            timeout: 3000,
          });
        } finally {
          this.loading = false;
        }
      },
      deleteDiscount() {
        Notiflix.Confirm.show(
          'Delete Discount',
          'Are you sure you want to delete this discount?',
          'Yes, delete it',
          'Cancel',
          async () => {
            try {
              const response = await fetch(urlDiscountDelete, {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                  discount_id: DISCOUNT_ID
                })
              });
              const data = await response.json();
              if (data.status === 'success') {
                Notiflix.Notify.success('Discount deleted successfully', {
                  position: 'right-top',
                  timeout: 2000,
                });
                setTimeout(() => {
                  window.location.href = urlDiscounts;
                }, 1500);
              } else {
                Notiflix.Notify.failure(data.message || 'Failed to delete discount', {
                  position: 'right-top',
                  timeout: 3000,
                });
              }
            } catch (error) {
              console.error('Failed to delete discount:', error);
              Notiflix.Notify.failure('Failed to delete discount', {
                position: 'right-top',
                timeout: 3000,
              });
            }
          },
          () => {
            // Cancel callback
          }
        );
      },
      formatDate(dateString) {
        if (!dateString) return '-';
        const date = new Date(dateString);
        const day = date.getDate().toString().padStart(2, '0');
        const month = date.toLocaleString('en-GB', { month: 'short' });
        const year = date.getFullYear();
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${day} ${month} ${year}<br><small class="text-muted">${hours}:${minutes}</small>`;
      },
      getStatusBadgeClass(status) {
        switch (status) {
          case 'active':
            return 'bg-success';
          case 'inactive':
            return 'bg-secondary';
          case 'draft':
            return 'bg-warning';
          case 'expired':
            return 'bg-danger';
          default:
            return 'bg-secondary';
        }
      }
    }
  }).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initDiscountView();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initDiscountView();
      }
    });
  }
});
