# mage-clean-media

`mage-clean-media` is a Go-based tool inspired by the [CleanMedia](https://github.com/cap340/CleanMedia) project. It is designed to clean up product media files and database records related to products that no longer exist or are redundant in a Magento-based e-commerce platform.

Follows the same logic as [CleanMedia](https://github.com/cap340/CleanMedia), the goal of this project was to just make that process faster and learn some Go in the process.

## Features

- **Media Cleaning**: Detects and removes unused media files associated with deleted products.
- **Database Cleanup**: Removes database records related to non-existent or redundant products.
- **Reporting**: Generates report of the cleaning process for transparency and auditing.
- **No module**: Requires no Magento 2 module installation. GO produces a binary executable that an be run on any system.

### Example

```bash
-mage-root /var/www/heals/web/ -host 127.0.0.1:40000 -name dbname -user root -password root -dry-run
```

### Options

- **-dry-run**: Runs the script without deleting files or DB records. (default true)
- **-dummy-data**: Set flag to generate a set of dummy image data.
- **-host string**: Database host (required).
- **-image-count int**: Define the number of images to generate with dummy data option. (default 500)
- **-mage-root string**: Declare the absolute path to the root of your Magento installation.
- **-name string**: Database name (required).
- **-no-cache**: Exclude files from the catalog/product/cache directory. (default true)
- **-password string**: Database password (required).
- **-user string**: Database username (required).

## Generating Dummy Image Data

You can use the `-dummy-data` flag to generate a set of dummy image data.

## Database Configuration

Ensure that the database credentials provided in the command options are correct. The tool requires access to the Magento database to identify and clean up records related to non-existent or redundant products.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Acknowledgments

Heavily inspired by [CleanMedia](https://github.com/cap340/CleanMedia)