function initCategoryManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.categoryManagerUrls || {},
        categories: [],
        selectedCategories: [],
        selectAll: false,
        loading: false,
        error: null,
        showFilterModal: false,
        currentPage: 0,
        perPage: 10,
        totalCategories: 0,
        filters: {
          status: ''
        },
        sortByColumn: 'created_at',
        sortOrder: 'desc'
      };
    },
    computed: {
      totalPages() {
        return Math.ceil(this.totalCategories / this.perPage);
      },
      filterStatus() {
        const parts = [];
        if (this.filters.status) parts.push(`status: ${this.filters.status}`);

        if (parts.length === 0) return 'Showing all categories';
        return 'Showing categories with ' + parts.join(', ');
      },
      hasActiveFilters() {
        return this.filters.status !== '';
      }
    },
    mounted() {
      const urlParams = new URLSearchParams(window.location.search);
      this.filters.status = urlParams.get('status') || '';
      this.sortByColumn = urlParams.get('sort_by') || 'created_at';
      this.sortOrder = urlParams.get('sort_order') || 'desc';
      this.currentPage = parseInt(urlParams.get('page') || '0', 10);
      this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
      this.loadCategories();
    },
    methods: {
      async loadCategories() {
        this.loading = true;
        this.error = null;
        try {
          const response = await fetch(this.urls.loadCategories, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              page: this.currentPage,
              per_page: this.perPage,
              status: this.filters.status,
              sort_by: this.sortByColumn,
              sort: this.sortOrder
            })
          });
          const result = await response.json();
          if (result.status === 'success') {
            if (result.data.urls) {
              this.urls = { ...this.urls, ...result.data.urls };
            }
            this.categories = result.data.categories;
            this.totalCategories = result.data.total || 0;
          } else {
            this.error = result.message || 'Failed to load categories';
          }
        } catch (error) {
          console.error('Failed to load categories:', error);
          this.error = 'Failed to load categories';
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
        const maxPage = Math.ceil(this.totalCategories / this.perPage) - 1;
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
        params.set('page', this.currentPage);
        params.set('per_page', this.perPage);
        params.set('sort_order', this.sortOrder);
        params.set('sort_by', this.sortByColumn);

        const newUrl = `${window.location.pathname}?${params.toString()}`;
        window.history.pushState({}, '', newUrl);

        this.closeFilterModal();
        this.loadCategories();
      },
      clearFilters() {
        this.filters = {
          status: ''
        };
        this.currentPage = 0;
        this.applyFilters();
      },
      async deleteCategory(categoryId) {
        Swal.fire({
          title: 'Delete Category',
          text: 'Are you sure you want to delete this category?',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete it'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.categoryDelete, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                category_id: categoryId
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.categories = this.categories.filter(c => c.id !== categoryId);
              Swal.fire('Success', 'Category deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete category', 'error');
            }
          } catch (error) {
            console.error('Failed to delete category:', error);
            Swal.fire('Error', 'Failed to delete category', 'error');
          }
        });
      },
      async deleteSelectedCategories() {
        Swal.fire({
          title: 'Delete Categories',
          text: `Are you sure you want to delete ${this.selectedCategories.length} category(ies)?`,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonColor: '#d33',
          cancelButtonColor: '#3085d6',
          confirmButtonText: 'Yes, delete them'
        }).then(async (result) => {
          if (!result.isConfirmed) return;
          try {
            const response = await fetch(this.urls.categoryDeleteSelected, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                bulk_category_ids: this.selectedCategories
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.categories = this.categories.filter(c => !this.selectedCategories.includes(c.id));
              this.selectedCategories = [];
              this.selectAll = false;
              Swal.fire('Success', 'Categories deleted successfully', 'success');
            } else {
              Swal.fire('Error', result.message || 'Failed to delete categories', 'error');
            }
          } catch (error) {
            console.error('Failed to delete categories:', error);
            Swal.fire('Error', 'Failed to delete categories', 'error');
          }
        });
      },
      toggleSelectAll() {
        if (this.selectAll) {
          this.selectedCategories = this.categories.map(c => c.id);
        } else {
          this.selectedCategories = [];
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
    initCategoryManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initCategoryManager();
      }
    });
  }
});
