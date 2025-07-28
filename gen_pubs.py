import os
import frontmatter
from pybtex.database import parse_file

# --- (El código para construir el mapa de ORCID sigue igual) ---
print("🔍 Escaneando perfiles de miembros para crear mapa de ORCID...")
orcid_to_member_map = {}
# ... (código existente aquí) ...
print("\n--- Mapa de ORCID creado con éxito ---\n")


# --- PROCESAR EL ARCHIVO BIBTEX (CON LÓGICA MEJORADA) ---
bib_file = 'referencias.bib'
data_dir = 'data/publications/'
content_dir = 'content/publications/'

os.makedirs(data_dir, exist_ok=True)
os.makedirs(content_dir, exist_ok=True)

bib_data = parse_file(bib_file)
print(f"📚 Procesando {len(bib_data.entries)} publicaciones desde {bib_file}...")

for key, entry in bib_data.entries.items():
    fields = entry.fields
    
    # --- LÓGICA DE DETECCIÓN DE MIEMBROS POR ORCID ---
    orcids_in_pub = fields.get('note', '').replace('ORCID:', '').split(',')
    found_members = []
    for orcid in orcids_in_pub:
        orcid = orcid.strip()
        if orcid in orcid_to_member_map:
            member_slug = orcid_to_member_map[orcid]
            if member_slug not in found_members:
                found_members.append(member_slug)
    
    # --- ¡NUEVO! OBTENER LA LISTA COMPLETA DE AUTORES ---
    # Esto extrae todos los nombres de los autores como texto.
    all_authors_list = [str(person) for person in entry.persons.get('author', [])]

    # Crear el diccionario de datos para el YAML
    yaml_data = {
        'title': fields.get('title', 'No Title').replace('{', '').replace('}', ''),
        'abstract': fields.get('abstract', '').replace('{', '').replace('}', ''),
        'journal': fields.get('journal', ''),
        'year': fields.get('year', ''),
        'doi': fields.get('doi', ''),
        # Lista de todos los autores (para mostrar)
        'all_authors': all_authors_list,
        # Lista de solo los miembros internos (para enlazar)
        'members': found_members,
    }
    
    # Escribir el archivo .yaml
    yaml_filename = os.path.join(data_dir, f"{key}.yaml")
    with open(yaml_filename, 'w', encoding='utf-8') as f:
        # Usaremos un formato más robusto para listas
        import yaml
        yaml.dump(yaml_data, f, allow_unicode=True)
            
    # Crear el archivo puntero .md
    md_filename = os.path.join(content_dir, f"{key}.md")
    with open(md_filename, 'w', encoding='utf-8') as f:
        # Pasamos los miembros aquí también para que Hugo los reconozca como taxonomía
        f.write(f'+++\ndatafile: "publications/{key}.yaml"\nmembers: {found_members}\n+++')

print("\n✅ ¡Proceso completado!")