function initProductView() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        loading: true,
        product: null
      };
    },
    mounted() {
      this.loadProductDetails();
    },
    methods: {
      async loadProductDetails() {
        this.loading = true;
        try {
          const response = await fetch(urlLoadDetails, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              product_id: PRODUCT_ID
            })
          });
          const data = await response.json();

          if (data.status === 'success') {
            this.product = data.data || null;
          } else {
            Notiflix.Notify.failure(data.message || 'Failed to load product details', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        } catch (error) {
          console.error('Error loading product details:', error);
          Notiflix.Notify.failure(error.message || 'Failed to load product details', {
            position: 'right-top',
            timeout: 3000,
          });
        } finally {
          this.loading = false;
        }
      },
      deleteProduct() {
        Notiflix.Confirm.show(
          'Delete Product',
          'Are you sure you want to delete this product?',
          'Yes, delete it',
          'Cancel',
          async () => {
            try {
              const response = await fetch(urlProductDelete, {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                  product_id: PRODUCT_ID
                })
              });
              const data = await response.json();
              if (data.status === 'success') {
                Notiflix.Notify.success('Product deleted successfully', {
                  position: 'right-top',
                  timeout: 2000,
                });
                setTimeout(() => {
                  window.location.href = urlProducts;
                }, 1500);
              } else {
                Notiflix.Notify.failure(data.message || 'Failed to delete product', {
                  position: 'right-top',
                  timeout: 3000,
                });
              }
            } catch (error) {
              console.error('Failed to delete product:', error);
              Notiflix.Notify.failure('Failed to delete product', {
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
          default:
            return 'bg-secondary';
        }
      }
    }
  }).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initProductView();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initProductView();
      }
    });
  }
});
