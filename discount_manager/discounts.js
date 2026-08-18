function initDiscountManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.discountManagerUrls || {},
        discounts: [],
        selectedDiscounts: [],
        selectAll: false,
        loading: false,
        error: null,
        showFilterModal: false,
        showCreateModal: false,
        creating: false,
        newDiscount: {
          title: ''
        },
        currentPage: 0,
        perPage: 10,
        totalDiscounts: 0,
        filters: {
          status: '',
          created_from: '',
          created_to: ''
        },
        sortByColumn: 'created_at',
        sortOrder: 'desc'
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalDiscounts / this.perPage);
      },
      filterStatus() {
        const parts = [];
        if (this.filters.status) parts.push(`status: ${this.filters.status}`);
        if (this.filters.created_from) parts.push(`from: ${this.filters.created_from}`);
        if (this.filters.created_to) parts.push(`to: ${this.filters.created_to}`);

        if (parts.length === 0) return 'Showing all discounts';
        return 'Showing discounts with ' + parts.join(', ');
      },
      hasActiveFilters() {
        return this.filters.status !== '' ||
               this.filters.created_from !== '' ||
               this.filters.created_to !== '';
      }
    },
    mounted() {
      const urlParams = new URLSearchParams(window.location.search);
      this.filters.status = urlParams.get('status') || '';
      this.filters.created_from = urlParams.get('created_from') || '';
      this.filters.created_to = urlParams.get('created_to') || '';
      this.sortByColumn = urlParams.get('sort_by') || 'created_at';
      this.sortOrder = urlParams.get('sort_order') || 'desc';
      this.currentPage = parseInt(urlParams.get('page') || '0', 10);
      this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
      this.loadDiscounts();
    },
    methods: {
      async loadDiscounts() {
        this.loading = true;
        this.error = null;
        try {
          const response = await fetch(this.urls.loadDiscounts, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              page: this.currentPage,
              per_page: this.perPage,
              status: this.filters.status,
              created_from: this.filters.created_from,
              created_to: this.filters.created_to,
              sort_by: this.sortByColumn,
              sort: this.sortOrder
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            if (result.data.urls) {
              this.urls = { ...this.urls, ...result.data.urls };
            }
            this.discounts = result.data.discounts;
            this.totalDiscounts = result.data.total || 0;
          } else {
            this.error = result.message || 'Failed to load discounts';
          }
        } catch (error) {
          console.error('Failed to load discounts:', error);
          this.error = 'Failed to load discounts';
        } finally {
          this.loading = false;
        }
      },
      sortBy(column) {
        if (this.sortByColumn === column) {
          this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
        } else {
          this.sortByColumn = column;
          this.sortOrder = 'asc';
        }
        this.currentPage = 0;
        this.applyFilters();
      },
      goToPage(page) {
        if (page < 0) return;
        const maxPage = Math.ceil(this.totalDiscounts / this.perPage) - 1;
        if (page > maxPage) page = maxPage;
        this.currentPage = page;
        this.applyFilters();
      },
      changePerPage() {
        this.perPage = parseInt(this.perPage, 10);
        this.currentPage = 0;
        this.applyFilters();
      },
      openFilterModal() {
        this.showFilterModal = true;
      },
      closeFilterModal() {
        this.showFilterModal = false;
      },
      applyFilters() {
        const params = new URLSearchParams();
        if (this.filters.status) params.set('status', this.filters.status);
        if (this.filters.created_from) params.set('created_from', this.filters.created_from);
        if (this.filters.created_to) params.set('created_to', this.filters.created_to);
        params.set('page', this.currentPage);
        params.set('per_page', this.perPage);
        params.set('sort_order', this.sortOrder);
        params.set('sort_by', this.sortByColumn);

        const newUrl = `${window.location.pathname}?${params.toString()}`;
        window.history.pushState({}, '', newUrl);

        this.closeFilterModal();
        this.loadDiscounts();
      },
      clearFilters() {
        this.filters = {
          status: '',
          created_from: '',
          created_to: ''
        };
        this.currentPage = 0;
        this.applyFilters();
      },
      async createDiscount() {
        if (!this.newDiscount.title.trim()) {
          Notiflix.Notify.failure('Title is required', {
            position: 'right-top',
            timeout: 3000,
          });
          return;
        }

        this.creating = true;
        try {
          const response = await fetch(this.urls.createDiscount, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              title: this.newDiscount.title.trim()
            })
          });
          const result = await response.json();

          if (result.status === 'success') {
            Notiflix.Notify.success('Discount created successfully', {
              position: 'right-top',
              timeout: 3000,
            });
            this.showCreateModal = false;
            this.newDiscount.title = '';
            this.loadDiscounts();
          } else {
            Notiflix.Notify.failure(result.message || 'Failed to create discount', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        } catch (error) {
          console.error('Failed to create discount:', error);
          Notiflix.Notify.failure('Failed to create discount', {
            position: 'right-top',
            timeout: 3000,
          });
        } finally {
          this.creating = false;
        }
      },
      async deleteDiscount(discountId) {
        Swal.fire({
          title: 'Delete Discount',
          text: 'Are you sure you want to delete this discount?',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete it'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.deleteDiscount, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                discount_id: discountId
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.discounts = this.discounts.filter(d => d.id !== discountId);
              Swal.fire('Success', 'Discount deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete discount', 'error');
            }
          } catch (error) {
            console.error('Failed to delete discount:', error);
            Swal.fire('Error', 'Failed to delete discount', 'error');
          }
        });
      },
      async deleteSelectedDiscounts() {
        Swal.fire({
          title: 'Delete Discounts',
          text: `Are you sure you want to delete ${this.selectedDiscounts.length} discount(s)?`,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete them'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.deleteSelectedDiscounts, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                bulk_discount_ids: this.selectedDiscounts
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.discounts = this.discounts.filter(d => !this.selectedDiscounts.includes(d.id));
              this.selectedDiscounts = [];
              this.selectAll = false;
              Swal.fire('Success', 'Discounts deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete discounts', 'error');
            }
          } catch (error) {
            console.error('Failed to delete discounts:', error);
            Swal.fire('Error', 'Failed to delete discounts', 'error');
          }
        });
      },
      toggleSelectAll() {
        if (this.selectAll) {
          this.selectedDiscounts = this.discounts.map(d => d.id);
        } else {
          this.selectedDiscounts = [];
        }
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
      },
      formatDate(dateString) {
        if (!dateString) return '';
        const date = new Date(dateString);
        return date.toLocaleDateString();
      }
    }
  }).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initDiscountManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initDiscountManager();
      }
    });
  }
});
