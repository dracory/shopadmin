function initCategoriesApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initCategoriesApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                collapsed: false,
                categories: [],
                selectedCategoryIDs: [],
                originalCategoryIDs: [],
                productID: categoriesProductId
            };
        },
        computed: {
            hasChanges() {
                if (this.selectedCategoryIDs.length !== this.originalCategoryIDs.length) {
                    return true;
                }
                const sortedSelected = [...this.selectedCategoryIDs].sort();
                const sortedOriginal = [...this.originalCategoryIDs].sort();
                for (let i = 0; i < sortedSelected.length; i++) {
                    if (sortedSelected[i] !== sortedOriginal[i]) {
                        return true;
                    }
                }
                return false;
            }
        },
        methods: {
            async loadCategories() {
                this.loading = true;
                try {
                    const response = await fetch(categoriesUrlLoad, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ product_id: this.productID })
                    });

                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        this.categories = data.categories || [];
                        this.selectedCategoryIDs = data.selected || [];
                        this.originalCategoryIDs = [...this.selectedCategoryIDs];
                    } else {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    console.error('Error loading categories:', error);
                    Notiflix.Notify.failure('Failed to load categories: ' + error.message);
                } finally {
                    this.loading = false;
                }
            },

            async saveCategories() {
                this.loading = true;
                try {
                    const response = await fetch(categoriesUrlSave, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            product_id: this.productID,
                            category_ids: this.selectedCategoryIDs
                        })
                    });

                    const data = await response.json();
                    if (data.status === 'success') {
                        Notiflix.Notify.success(data.message || 'Categories saved successfully');
                        this.originalCategoryIDs = [...this.selectedCategoryIDs];
                    } else {
                        Notiflix.Notify.failure(data.message || 'Failed to save categories');
                    }
                } catch (error) {
                    console.error('Error saving categories:', error);
                    Notiflix.Notify.failure('Failed to save categories: ' + error.message);
                } finally {
                    this.loading = false;
                }
            }
        },
        mounted() {
            this.loadCategories();
        }
    }).mount('#categories-app');
}

document.addEventListener('DOMContentLoaded', initCategoriesApp);
