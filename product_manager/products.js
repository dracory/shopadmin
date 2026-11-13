function initProductManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.productManagerUrls || {},
        products: [],
        selectedProducts: [],
        selectAll: false,
        loading: false,
        error: null,
        showFilterModal: false,
        showCreateModal: false,
        creating: false,
        newProduct: {
          title: ''
        },
        currentPage: 0,
        perPage: 10,
        totalProducts: 0,
        // Applied conditions (sent to the server)
        conditions: [],
        // Modal state — cloned from conditions when opening
        modalConditions: [],
        sortByColumn: 'created_at',
        sortOrder: 'desc',
        // Filter condition options — must match the Go handler fields.
        conditionOptions: [
          { value: 'search', label: 'Search', operators: ['contains'], inputType: 'text', placeholder: 'title...' },
          { value: 'status', label: 'Status', operators: ['equals'], inputType: 'select', placeholder: '', options: [
            { value: 'active', label: 'Active' },
            { value: 'inactive', label: 'Inactive' },
            { value: 'draft', label: 'Draft' },
          ] },
          { value: 'id', label: 'ID', operators: ['equals'], inputType: 'text', placeholder: 'product id' },
          { value: 'category_id', label: 'Category', operators: ['equals'], inputType: 'select', placeholder: '', options: (window.productManagerCategories || []).map(c => ({ value: c.id, label: c.title })) },
          { value: 'date_from', label: 'Created after', operators: ['equals'], inputType: 'date', placeholder: '' },
          { value: 'date_to', label: 'Created before', operators: ['equals'], inputType: 'date', placeholder: '' },
        ],
        opLabels: { equals: '=', contains: 'contains' }
      };
    },
  computed: {
    totalPages() {
      return Math.ceil(this.totalProducts / this.perPage);
    },
    filterStatus() {
      if (this.conditions.length === 0) return 'Showing all products';
      return this.conditions.length + ' filter(s) active';
    },
    hasActiveFilters() {
      return this.conditions.length > 0;
    }
  },
  mounted() {
    const urlParams = new URLSearchParams(window.location.search);

    // Try JSON filters first (new format)
    const filtersRaw = urlParams.get('filters');
    if (filtersRaw) {
      try {
        const parsed = JSON.parse(filtersRaw);
        if (Array.isArray(parsed)) {
          this.conditions = parsed.filter(c => c && c.field && c.value);
        }
      } catch (e) { /* ignore malformed */ }
    } else {
      // Fallback: legacy individual params
      if (urlParams.get('status')) this.conditions.push({ field: 'status', operator: 'equals', value: urlParams.get('status') });
      if (urlParams.get('created_from')) this.conditions.push({ field: 'date_from', operator: 'equals', value: urlParams.get('created_from') });
      if (urlParams.get('created_to')) this.conditions.push({ field: 'date_to', operator: 'equals', value: urlParams.get('created_to') });
    }

    this.sortByColumn = urlParams.get('sort_by') || 'created_at';
    this.sortOrder = urlParams.get('sort_order') || 'desc';
    this.currentPage = parseInt(urlParams.get('page') || '0', 10);
    this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
    this.loadProducts();
  },
  methods: {
    async loadProducts() {
      this.loading = true;
      this.error = null;
      try {
        const response = await fetch(this.urls.loadProducts, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            page: this.currentPage,
            per_page: this.perPage,
            conditions: this.conditions,
            sort_by: this.sortByColumn,
            sort: this.sortOrder
          })
        });
        const result = await response.json();
        console.log('Products response:', result);
        if (result.status === 'success') {
          this.products = result.data.products;
          this.totalProducts = result.data.total || 0;
        } else {
          this.error = result.message || 'Failed to load products';
        }
      } catch (error) {
        console.error('Failed to load products:', error);
        this.error = 'Failed to load products';
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
      const maxPage = Math.ceil(this.totalProducts / this.perPage) - 1;
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
      this.modalConditions = this.conditions.map(c => ({ ...c }));
      if (this.modalConditions.length === 0) {
        this.modalConditions.push({ field: '', operator: 'equals', value: '' });
      }
      this.showFilterModal = true;
    },
    closeFilterModal() {
      this.showFilterModal = false;
    },
    addModalCondition() {
      this.modalConditions.push({ field: '', operator: 'equals', value: '' });
    },
    removeModalCondition(i) {
      this.modalConditions.splice(i, 1);
    },
    clearModalConditions() {
      this.modalConditions = [];
    },
    onFieldChange(c) {
      const opt = this.getFieldOpt(c.field);
      if (opt) {
        c.operator = opt.operators[0];
      }
      c.value = '';
    },
    applyConditions() {
      this.conditions = this.modalConditions.filter(c => c.field && String(c.value).trim());
      this.showFilterModal = false;
      this.currentPage = 0;
      this.applyFilters();
    },
    applyFilters() {
      const params = new URLSearchParams();
      const existing = new URLSearchParams(window.location.search);
      const controller = existing.get('controller');
      if (controller) params.set('controller', controller);

      if (this.conditions.length > 0) params.set('filters', JSON.stringify(this.conditions));
      params.set('page', this.currentPage);
      params.set('per_page', this.perPage);
      params.set('sort_order', this.sortOrder);
      params.set('sort_by', this.sortByColumn);

      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newUrl);

      this.loadProducts();
    },
    clearFilters() {
      this.conditions = [];
      this.currentPage = 0;
      this.applyFilters();
    },
    // === Filter option helpers ===
    getFieldOpt(v) { return this.conditionOptions.find(o => o.value === v); },
    opLabel(op) { return this.opLabels[op] || op; },
    getFieldOperators(v) { const o = this.getFieldOpt(v); return o ? o.operators : ['equals']; },
    getFieldInputType(v) { const o = this.getFieldOpt(v); return o ? o.inputType : 'text'; },
    getFieldPlaceholder(v) { const o = this.getFieldOpt(v); return o ? o.placeholder : ''; },
    getFieldSelectOptions(v) { const o = this.getFieldOpt(v); return o && o.options ? o.options : []; },
    sortIcon(column) {
      if (this.sortByColumn !== column) return 'bi bi-arrow-down-up text-muted';
      return this.sortOrder === 'asc' ? 'bi bi-arrow-up' : 'bi bi-arrow-down';
    },
    async deleteProduct(productId) {
      Notiflix.Confirm.show(
        'Delete Product',
        'Are you sure you want to delete this product?',
        'Yes, delete it',
        'Cancel',
        async () => {
          try {
            const response = await fetch(this.urls.productDelete, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                product_id: productId
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.products = this.products.filter(p => p.id !== productId);
              Notiflix.Notify.success('Product deleted successfully', {
                position: 'right-top',
                timeout: 3000,
              });
            } else {
              Notiflix.Notify.failure(result.message || 'Failed to delete product', {
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
    async deleteSelectedProducts() {
      Notiflix.Confirm.show(
        'Delete Products',
        `Are you sure you want to delete ${this.selectedProducts.length} product(s)?`,
        'Yes, delete them',
        'Cancel',
        async () => {
          try {
            const response = await fetch(this.urls.productDeleteSelected, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                bulk_product_ids: this.selectedProducts
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.products = this.products.filter(p => !this.selectedProducts.includes(p.id));
              this.selectedProducts = [];
              this.selectAll = false;
              Notiflix.Notify.success('Products deleted successfully', {
                position: 'right-top',
                timeout: 3000,
              });
            } else {
              Notiflix.Notify.failure(result.message || 'Failed to delete products', {
                position: 'right-top',
                timeout: 3000,
              });
            }
          } catch (error) {
            console.error('Failed to delete products:', error);
            Notiflix.Notify.failure('Failed to delete products', {
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
    toggleSelectAll() {
      if (this.selectAll) {
        this.selectedProducts = this.products.map(p => p.id);
      } else {
        this.selectedProducts = [];
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
    },
    async createProduct() {
      if (!this.newProduct.title.trim()) {
        Notiflix.Notify.failure('Title is required', {
          position: 'right-top',
          timeout: 3000,
        });
        return;
      }

      this.creating = true;
      try {
        const response = await fetch(this.urls.createProduct, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            title: this.newProduct.title.trim()
          })
        });
        const result = await response.json();

        if (result.status === 'success') {
          Notiflix.Notify.success('Product created successfully', {
            position: 'right-top',
            timeout: 3000,
          });
          this.showCreateModal = false;
          this.newProduct.title = '';
          this.loadProducts();
        } else {
          Notiflix.Notify.failure(result.message || 'Failed to create product', {
            position: 'right-top',
            timeout: 3000,
          });
        }
      } catch (error) {
        console.error('Failed to create product:', error);
        Notiflix.Notify.failure('Failed to create product', {
          position: 'right-top',
          timeout: 3000,
        });
      } finally {
        this.creating = false;
      }
    }
  }
}).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initProductManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initProductManager();
      }
    });
  }
});
