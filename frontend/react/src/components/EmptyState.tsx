interface EmptyStateProps {
    title: string;
    description: string;
    actionLabel?: string;
    onAction?: () => void;
}

export default function EmptyState({
    title,
    description,
    actionLabel,
    onAction,
}: EmptyStateProps) {
    return (
        <div className="empty-state">
            <div className="empty-state__icon">
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="1.5"
                    width="80"
                    height="80"
                >
                    <rect x="2" y="7" width="20" height="11" rx="3" />
                    <circle cx="7" cy="15" r="2" />
                    <circle cx="17" cy="15" r="2" />
                    <path d="M5 7l2-3h10l2 3" />
                </svg>
            </div>
            <h3 className="empty-state__title">{title}</h3>
            <p className="empty-state__text">{description}</p>
            {actionLabel && onAction && (
                <button className="btn btn--primary" onClick={onAction}>
                    {actionLabel}
                </button>
            )}
        </div>
    );
}
