interface PaginationProps {
    offset: number;
    limit: number;
    total: number;
    onPageChange: (newOffset: number) => void;
}

export default function Pagination({
    offset,
    limit,
    total,
    onPageChange,
}: PaginationProps) {
    const currentPage = Math.floor(offset / limit) + 1;
    const totalPages = Math.ceil(total / limit);

    if (totalPages <= 1) return null;

    const pages: number[] = [];
    const maxVisible = 5;
    let start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
    const end = Math.min(totalPages, start + maxVisible - 1);

    if (end - start + 1 < maxVisible) {
        start = Math.max(1, end - maxVisible + 1);
    }

    for (let i = start; i <= end; i++) {
        pages.push(i);
    }

    return (
        <div className="pagination">
            <button
                className="pagination__btn"
                disabled={currentPage === 1}
                onClick={() => onPageChange(0)}
                title="First page"
            >
                ««
            </button>
            <button
                className="pagination__btn"
                disabled={currentPage === 1}
                onClick={() => onPageChange(offset - limit)}
                title="Previous page"
            >
                «
            </button>

            {start > 1 && <span className="pagination__ellipsis">…</span>}

            {pages.map((page) => (
                <button
                    key={page}
                    className={`pagination__btn ${page === currentPage ? "pagination__btn--active" : ""}`}
                    onClick={() => onPageChange((page - 1) * limit)}
                >
                    {page}
                </button>
            ))}

            {end < totalPages && <span className="pagination__ellipsis">…</span>}

            <button
                className="pagination__btn"
                disabled={currentPage === totalPages}
                onClick={() => onPageChange(offset + limit)}
                title="Next page"
            >
                »
            </button>
            <button
                className="pagination__btn"
                disabled={currentPage === totalPages}
                onClick={() => onPageChange((totalPages - 1) * limit)}
                title="Last page"
            >
                »»
            </button>

            <span className="pagination__info">
                {offset + 1}–{Math.min(offset + limit, total)} of {total}
            </span>
        </div>
    );
}
