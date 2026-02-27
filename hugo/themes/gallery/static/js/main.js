document.addEventListener('DOMContentLoaded', function() {
    const lightbox = document.getElementById('lightbox');
    const lightboxImg = lightbox.querySelector('.lightbox-image');
    const lightboxInfo = lightbox.querySelector('.lightbox-info');
    const closeBtn = lightbox.querySelector('.lightbox-close');
    const prevBtn = lightbox.querySelector('.lightbox-prev');
    const nextBtn = lightbox.querySelector('.lightbox-next');
    
    let currentIndex = 0;
    let photos = [];
    let currentTaxonomy = '';
    let currentTerm = '';
    
    function getPhotoCards() {
        return document.querySelectorAll('.photo-card[data-permalink]');
    }
    
    function buildPhotosFromContext(card) {
        const taxonomy = card.dataset.taxonomy;
        const term = card.dataset.term;
        
        if (!taxonomy || !term) {
            return Array.from(getPhotoCards()).map(c => {
                const fullSrc = c.querySelector('img')?.dataset?.fullSrc;
                return {
                    src: fullSrc || c.querySelector('img').src,
                    link: c.dataset.permalink
                };
            });
        }
        
        const allCards = getPhotoCards();
        const contextCards = Array.from(allCards).filter(c => 
            c.dataset.taxonomy === taxonomy && c.dataset.term === term
        );
        
        return contextCards.map(c => {
            const fullSrc = c.querySelector('img')?.dataset?.fullSrc;
            return {
                src: fullSrc || c.querySelector('img').src,
                link: c.dataset.permalink
            };
        });
    }
    
    function getCurrentPhotoIndex(cards, clickedPermalink) {
        for (let i = 0; i < cards.length; i++) {
            if (cards[i].link === clickedPermalink) return i;
        }
        return 0;
    }
    
    function openLightbox(src, link, taxonomy, term) {
        lightboxImg.src = src;
        lightbox.classList.add('active');
        document.body.style.overflow = 'hidden';
        
        if (link && taxonomy && term) {
            lightboxInfo.innerHTML = `<a href="${link}?from=${encodeURIComponent(taxonomy)}:${encodeURIComponent(term)}" class="lightbox-link">View Photo Details</a>`;
        } else if (link) {
            lightboxInfo.innerHTML = `<a href="${link}" class="lightbox-link">View Photo Details</a>`;
        } else {
            lightboxInfo.innerHTML = '';
        }
    }
    
    function closeLightbox() {
        lightbox.classList.remove('active');
        document.body.style.overflow = '';
    }
    
    function showPrev() {
        currentIndex = (currentIndex - 1 + photos.length) % photos.length;
        lightboxImg.src = photos[currentIndex].src;
        if (photos[currentIndex].link && currentTaxonomy && currentTerm) {
            lightboxInfo.innerHTML = `<a href="${photos[currentIndex].link}?from=${encodeURIComponent(currentTaxonomy)}:${encodeURIComponent(currentTerm)}" class="lightbox-link">View Photo Details</a>`;
        } else if (photos[currentIndex].link) {
            lightboxInfo.innerHTML = `<a href="${photos[currentIndex].link}" class="lightbox-link">View Photo Details</a>`;
        } else {
            lightboxInfo.innerHTML = '';
        }
    }
    
    function showNext() {
        currentIndex = (currentIndex + 1) % photos.length;
        lightboxImg.src = photos[currentIndex].src;
        if (photos[currentIndex].link && currentTaxonomy && currentTerm) {
            lightboxInfo.innerHTML = `<a href="${photos[currentIndex].link}?from=${encodeURIComponent(currentTaxonomy)}:${encodeURIComponent(currentTerm)}" class="lightbox-link">View Photo Details</a>`;
        } else if (photos[currentIndex].link) {
            lightboxInfo.innerHTML = `<a href="${photos[currentIndex].link}" class="lightbox-link">View Photo Details</a>`;
        } else {
            lightboxInfo.innerHTML = '';
        }
    }
    
    closeBtn.addEventListener('click', closeLightbox);
    prevBtn.addEventListener('click', showPrev);
    nextBtn.addEventListener('click', showNext);
    
    lightbox.addEventListener('click', function(e) {
        if (e.target === lightbox) {
            closeLightbox();
        }
    });
    
    document.addEventListener('keydown', function(e) {
        if (!lightbox.classList.contains('active')) return;
        
        if (e.key === 'Escape') closeLightbox();
        if (e.key === 'ArrowLeft') showPrev();
        if (e.key === 'ArrowRight') showNext();
    });
    
    // Touch/swipe support for lightbox
    let touchStartX = 0;
    lightbox.addEventListener('touchstart', e => {
        touchStartX = e.changedTouches[0].screenX;
    }, { passive: true });
    
    lightbox.addEventListener('touchend', e => {
        const touchEndX = e.changedTouches[0].screenX;
        const diff = touchStartX - touchEndX;
        if (Math.abs(diff) > 50) {
            if (diff > 0) showNext();
            else showPrev();
        }
    });
    
    const cards = getPhotoCards();
    
    cards.forEach(card => {
        const slideshowIcon = card.querySelector('.photo-slideshow-icon');
        
        if (slideshowIcon) {
            slideshowIcon.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                
                currentTaxonomy = card.dataset.taxonomy || '';
                currentTerm = card.dataset.term || '';
                
                photos = buildPhotosFromContext(card);
                currentIndex = getCurrentPhotoIndex(photos, card.dataset.permalink);
                
                const clickedLink = card.dataset.permalink;
                openLightbox(photos[currentIndex].src, clickedLink, currentTaxonomy, currentTerm);
            });
        }
    });
    
    document.querySelectorAll('.slideshow-btn').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            
            const taxonomy = this.dataset.taxonomy || '';
            const term = this.dataset.term || '';
            
            const grid = this.closest('.album-header, .page-header').nextElementSibling;
            const gridCards = grid.querySelectorAll('.photo-card');
            
            if (gridCards.length === 0) return;
            
            currentTaxonomy = taxonomy;
            currentTerm = term;
            
            photos = Array.from(gridCards).map(c => {
                const fullSrc = c.querySelector('img')?.dataset?.fullSrc;
                return {
                    src: fullSrc || c.querySelector('img').src,
                    link: c.dataset.permalink
                };
            });
            
            currentIndex = 0;
            const firstCard = gridCards[0];
            const firstLink = firstCard.dataset.permalink;
            openLightbox(photos[0].src, firstLink, currentTaxonomy, currentTerm);
        });
    });
});
