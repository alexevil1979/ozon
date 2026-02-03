-- Catalog: product attribute values (many-to-many product <-> attribute).
CREATE TABLE IF NOT EXISTS product_attribute_values (
    product_id   UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    attribute_id UUID NOT NULL REFERENCES attributes (id) ON DELETE CASCADE,
    value        TEXT NOT NULL,
    PRIMARY KEY (product_id, attribute_id)
);

CREATE INDEX IF NOT EXISTS idx_product_attribute_values_attribute_id ON product_attribute_values (attribute_id);
