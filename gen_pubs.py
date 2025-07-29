import os
import requests
import yaml # Lo mantenemos para leer los YAML existentes y evitar errores
import toml # La nueva librería para escribir TOML
import datetime
from slugify import slugify

# --- CONFIGURACIÓN ---
# Directorio donde se guardarán las publicaciones.
PUBLICATIONS_DIR = "content/publications"
# Tu email, es una buena práctica para usar APIs públicas.
CROSSREF_MAILTO = "dperez42@uasd.edu.do" 

def get_existing_dois(directory):
    """
    Escanea el directorio de publicaciones y extrae todos los DOI existentes,
    soportando tanto front matter YAML (---) como TOML (+++).
    """
    existing_dois = set()
    if not os.path.exists(directory):
        return existing_dois

    for filename in os.listdir(directory):
        if filename.endswith(".md"):
            filepath = os.path.join(directory, filename)
            try:
                with open(filepath, 'r', encoding='utf-8') as f:
                    content = f.read()
                    
                    # Determinar el tipo de front matter por el delimitador
                    if content.strip().startswith('---'):
                        delimiter = '---'
                        loader = yaml.safe_load
                    elif content.strip().startswith('+++'):
                        delimiter = '+++'
                        loader = toml.loads
                    else:
                        continue # No es un archivo con front matter que reconozcamos

                    parts = content.split(delimiter)
                    if len(parts) > 2:
                        front_matter = loader(parts[1])
                        if front_matter and 'doi' in front_matter:
                            existing_dois.add(str(front_matter['doi']).lower())

            except Exception as e:
                print(f"Advertencia: No se pudo leer el archivo {filepath}. Error: {e}")
    
    return existing_dois

def fetch_doi_metadata(doi):
    """
    Obtiene los metadatos de un DOI usando la API de Crossref.
    """
    print(f"Buscando metadatos para DOI: {doi}...")
    url = f"https://api.crossref.org/works/{doi}"
    headers = {
        'User-Agent': f'HugoPublicationsScript/1.0 (mailto:{CROSSREF_MAILTO})'
    }
    try:
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        print("Metadatos encontrados.")
        return response.json()['message']
    except requests.exceptions.RequestException as e:
        print(f"Error al obtener datos para el DOI {doi}: {e}")
        return None

def create_publication_file(metadata, staff):
    """
    Crea el archivo .md con el front matter en formato TOML.
    """
    # 1. Extraer y formatear datos
    title = metadata.get('title', ['Sin Título'])[0]
    
    try:
        date_parts = metadata.get('published', {}).get('date-parts', [[None]])[0]
        # TOML requiere un objeto de fecha y hora completo
        publish_date = datetime.datetime(date_parts[0], date_parts[1], date_parts[2])
    except (TypeError, IndexError, ValueError):
        publish_date = datetime.datetime.now()

    authors_list = []
    if 'author' in metadata:
        for author in metadata['author']:
            author_name = f"{author.get('given', '')} {author.get('family', '')}".strip()
            if author_name:
                authors_list.append(author_name)

    # 2. Construir el diccionario para el front matter
    front_matter = {
        'title': title,
        'date': publish_date, # Usamos el objeto datetime completo
        'doi': str(metadata.get('DOI', '')).lower(),
        'staff': staff,
        'publication': metadata.get('container-title', [''])[0],
        'authors': authors_list,
        'draft': False
    }

    # 3. Crear nombre de archivo
    first_author_family = metadata.get('author', [{}])[0].get('family', 'unknown')
    year = publish_date.strftime("%Y")
    slug_title = slugify(title, max_length=40)
    filename = f"{first_author_family.lower()}-{year}-{slug_title}.md"
    filepath = os.path.join(PUBLICATIONS_DIR, filename)

    # 4. Escribir el archivo con formato TOML
    print(f"Creando archivo: {filepath}")
    os.makedirs(PUBLICATIONS_DIR, exist_ok=True)
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write('+++\n')
        # Convertimos el diccionario a una cadena en formato TOML
        toml_string = toml.dumps(front_matter)
        f.write(toml_string)
        f.write('+++\n\n')
        f.write("\n")
    print("¡Archivo creado con éxito!")

def main():
    """
    Función principal del script.
    """
    # --- LISTA DE PUBLICACIONES PARA AÑADIR ---
    # Edita esta lista con los DOI y los miembros del instituto asociados.
    publications_to_add = [
        {"doi": "10.56048/MQR20225.9.2.2025.e664", "staff": ["franmis-rodriguez", "erika-montero"]},
        {"doi": "10.21071/edmetic.v10i2.13240", "staff": ["franmis-rodriguez"]},
        {"doi": "10.15517/revedu.v48i1.55892", "staff": ["franmis-rodriguez"]}, # Ejemplo sin miembros
        # Añade aquí más publicaciones
    ]
    # ---------------------------------------------

    print("--- Iniciando script para añadir publicaciones (formato TOML) ---")
    existing_dois = get_existing_dois(PUBLICATIONS_DIR)
    print(f"Se encontraron {len(existing_dois)} DOI existentes.")

    for pub in publications_to_add:
        doi_to_check = pub['doi'].lower()
        
        if doi_to_check in existing_dois:
            print(f"El DOI {doi_to_check} ya existe. Saltando.")
            continue
        
        metadata = fetch_doi_metadata(pub['doi'])
        
        if metadata:
            create_publication_file(metadata, pub['staff'])
        
        print("-" * 20)

    print("--- Script finalizado ---")

if __name__ == "__main__":
    main()